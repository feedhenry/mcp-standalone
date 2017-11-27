'use strict';

/**
 * @ngdoc component
 * @name mcp.component:mp-ios-setup
 * @description
 * # mp-ios-setup
 */
angular.module('mobileControlPanelApp').component('mpIosSetup', {
  template: `<div ng-if="$ctrl.app.clientType === 'iOS'">
              <h4>Installation</h4>
              <p>Add the following to your <code>Podfile</code></p>
<mp-prettify type="'bash'" code-class="'prettyprint'">
pod 'MobileCore', '~> 0.0.1'
</mp-prettify>

              <h4>Configuration</h4>
              <p>Add the following to a file named <code>mobile.plist</code> under your resources directory</p>
<mp-prettify type="'xml'" code-class="'prettyprint'">
&lt;?xml version="1.0" encoding="UTF-8"?&gt;
&lt;!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd"&gt;
&lt;dict&gt;
  &lt;key&gt;host&lt;/key&gt;
  &lt;string&gt;{{$ctrl.route}}&lt;/string&gt;
  &lt;key&gt;appID&lt;/key&gt;
  &lt;string&gt;{{$ctrl.app.id}}&lt;/string&gt;
  &lt;key&gt;apiKey&lt;/key&gt;
  &lt;string&gt;{{$ctrl.app.apiKey}}&lt;/string&gt;
&lt;/dict&gt;
&lt;/plist&gt;
</mp-prettify>

              <h4>SDK Initialisation</h4>
<mp-prettify type="'swift'" code-class="'prettyprint'">
var myDict: NSDictionary?
if let path = Bundle.main.path(forResource: "Config", ofType: "plist") {
  myDict = NSDictionary(contentsOfFile: path)
}
if let dict = myDict {
  // Use your dict here
}
</mp-prettify>

              <h3>Example Apps</h3>
              <div>
                <table class="table">
                  <thead>
                    <tr>
                      <th>Template Name</th>
                      <th>Description</th>
                      <th>Source</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr>
                      <td>Hello World</td>
                      <td>iOS swift hello world starter app. Shows you how to plugin the core sdk.</td>
                      <td><a href="https://github.com/feedhenry-templates/sync-ios-app">https://github.com/feedhenry-templates/ios-app</a></td>
                    </tr>
                    <tr>
                      <td>Sync iOS quick start</td>
                      <td>A starting point for building out an application that syncs data automatically to the cloud</td>
                      <td><a href="https://github.com/feedhenry-templates/sync-ios-app">https://github.com/feedhenry-templates/sync-ios-app</a></td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>`,
  bindings: {
    app: '<'
  }
});
