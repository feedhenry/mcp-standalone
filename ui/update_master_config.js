var yaml = require('js-yaml');
var fs = require('fs');

var mobileDir = '/var/lib/origin/openshift.local.config';
var mobileDistDir = mobileDir + '/dist';
var mobileViewsDir = mobileDir + '/public';
var mcpJSFiles = [mobileDistDir + '/mcp-vendor.js', mobileDistDir + '/mcp.js'];
var mcpCSSFiles = [mobileDistDir + '/mcp-vendor.css', mobileDistDir + '/mcp.css'];
var configFile = process.argv.slice(-1)[0];
var yamlFile = yaml.safeLoad(fs.readFileSync(configFile));


// Enable extension development
yamlFile.assetConfig.extensionDevelopment = true;

// Add mcp js files
yamlFile.assetConfig.extensionScripts = yamlFile.assetConfig.extensionScripts || [];
mcpJSFiles.forEach(function(mcpJSFile) {
  if (yamlFile.assetConfig.extensionScripts.indexOf(mcpJSFile) < 0) {
    yamlFile.assetConfig.extensionScripts.push(mcpJSFile);
  }
});

// Add mcp css files
yamlFile.assetConfig.extensionStylesheets = yamlFile.assetConfig.extensionStylesheets || [];
mcpCSSFiles.forEach(function(mcpCSSFile) {
  if (yamlFile.assetConfig.extensionStylesheets.indexOf(mcpCSSFile) < 0) {
    yamlFile.assetConfig.extensionStylesheets.push(mcpCSSFile);
  }
});

// Register mcp extension
yamlFile.assetConfig.extensions = yamlFile.assetConfig.extensions || [];
var mcpExtensionAdded = false;
yamlFile.assetConfig.extensions.forEach(function(extension) {
  if (extension.name === 'mcp') {
    mcpExtensionAdded = true;
  }
});

// TODO: use dist dir for views & directives too,
// and have a grunt task for copying those to dist?
if (!mcpExtensionAdded) {
  yamlFile.assetConfig.extensions.push({
    name: 'mcp',
    sourceDirectory: mobileViewsDir
  });
}

// write file
fs.writeFileSync(configFile, yaml.safeDump(yamlFile));