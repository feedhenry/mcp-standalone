
0.0.6 / 2017-11-02
==================

  * [make apbs script] updating Openshift template for APBs
  * Adding latest release tag
  * Merge pull request #192 from matzew/Deferred_Evaluation_for_release
  * :scream: nesting the check for modified files...
  * Merge pull request #193 from feedhenry/pass-through-secret-params-during-bind
  * pass through all params during a bind
  * Merge pull request #191 from finp/patch-1
  * update readme to add missing quote
  * Merge pull request #176 from matzew/Release_Process
  * Merge pull request #190 from feedhenry/fix-role-binding
  * fix issue with apply service account role binding
  * Merge pull request #188 from sedroche/service-integrations
  * Merge pull request #4 from aidenkeating/service-integration-fix
  * Remove repeated block of code
  * Merge pull request #187 from feedhenry/add-design-flats
  * Service integrations components
  * add ui flats doc
  * Merge pull request #166 from philbrookes/FH-4280
  * Merge pull request #186 from aidenkeating/unquote-asb-credentials
  * Unquote ASB credentials from base64
  * Add description to integrations and services and show them in the UI
  * Merge pull request #185 from aidenkeating/broker-image-tag
  * Use non-latest broker image tag
  * Merge pull request #183 from aidenkeating/remove-log-cli
  * Remove log from CLI to allow for piping
  * Merge pull request #182 from aidenkeating/update-asb-docker-credentials
  * Use clusterServiceClassExternalName
  * Update APBs
  * Use oc artifact from GitHub release
  * Update ASB Docker Credentials
  * Merge pull request #179 from aidenkeating/display-service-icons
  *  FH-4325 Display service icons
  * Merge pull request #178 from feedhenry/excluster-3scale-config
  * update domain to be more flavourful
  * Merge pull request #171 from sedroche/services-controller-refactor
  * Update keycloak metrics test with correct mock uri value
  * remove unneeded service from returned config
  * Fixup for origin 3.7 changes, keycloak host & sync graphs
  * Refactor charts into components
  * :lipstick: Section for MCP release steps and our independent APBs
  * :smiling_imp: No release with locally modified files
  * :space_invader: adding 'auto-commit' on CMD updates for APBs
  * Merge pull request #174 from aidenkeating/generic-client-config
  * Use default sync uri if 3scale integration isn't made
  * Merge pull request #175 from feedhenry/lock-down-origin-images
  * inset the headers
  * inset the headers
  * 🐛 Lock down origin images to have a stable local cluster
  * Implement GenericClientConfig
  * Merge pull request #149 from feedhenry/FH-4242-keycloak-bind
  * update logging
  * Implement GenericClientConfig
  * use bind api to enable integrations between services
  * Merge pull request #169 from feedhenry/maleck13-patch-1
  * Merge pull request #173 from matzew/No_Plastic
  * :lipstick: Using same image lnf
  * add useful bash function
  * Merge pull request #168 from matzew/Nice_badges
  * :lipstick: show some awesome badges
  * Merge pull request #167 from matzew/Release_Process_fixes
  * :whale: defaulting to docker.io registry as the DOCKERHOST
  * Merge pull request #164 from aidenkeating/3.7
  * Cleanup
  * Change serviceclass resource to clusterserviceclass
  * Reinclude TEMPLATE_VARS in ASB provision script
  * Exit on error for ansible service broker setup
  * Only 'become' sudo for linux
  * Download template to /tmp/
  * Use different dest paths for linux, and use sudo
  * Making specs compliant w/ the latest ASB
  * Allow linux artifact downloads and note temporary fix
  * 3.7

0.0.5 / 2017-10-24
==================

  * Merge pull request #161 from feedhenry/add-version
  * Merge branch 'master' into add-version
  * Merge pull request #165 from feedhenry/FH-4294-tags
  * fix ascii doc
  * Merge pull request #160 from feedhenry/FH-4294-tags
  * update README
  * add template version
  * add basic release doc
  * initial release doc and impl
  * allow a version to be set at build time
  * Merge pull request #159 from aidenkeating/FH-4207-keycloak-walkthrough
  * FH-4207 Rewrite some parts of walkthrough
  * FH-4207 Add Keycloak walkthrough
  * Merge pull request #157 from matzew/Remove_FH_sync_openshift_template
  * Merge pull request #156 from feedhenry/FH-4206-sync-and-keycloak-walkthrough
  * Add link to mobile cicd walkthrough
  * Address feedback
  * 📝 FH-4206 Add a walkthrough for Sync & Keycloak integration
  * Merge pull request #158 from aidenkeating/FH-4208-mobile-ci-cd-walkthrough
  * FH-4208 Walkthrough of Mobile CI/CD with MCP
  * Removing the OSCP template for fh-sync-server

0.0.4 / 2017-10-22
==================

  * initial release doc and impl
  * Merge pull request #151 from sedroche/tab-updates
  * Merge pull request #153 from finp/remove-sudo-from-walkthru
  * Merge pull request #155 from feedhenry/remove-master-build
  * remove master build
  * Merge pull request #152 from sedroche/modal-build-fix
  * remove sudo from ansible-galaxy command in walkthru
  * Annotate dependencies for minification
  * mp-tabs
  * Merge pull request #147 from philbrookes/update-local-dev-docs
  * skip syncing APB's by default
  * update docs for local APB development

0.0.3 / 2017-10-12
==================

  * Merge pull request #145 from feedhenry/metrics-enable-readme
  * Merge pull request #146 from feedhenry/use-origin36
  * fix metrics enabling via ansible env-var
  * readme updates for enabling metrics
  * use origin v3.6.0 instead of v3.6.0-rc.0
  * Merge pull request #144 from aidenkeating/FH-4245-rename-3scale-service
  * Expect type 3scale
  * Merge pull request #143 from aidenkeating/FH-4245-ensure-integration-appears-ui
  * Merge pull request #142 from philbrookes/FH-4258
  * Ensure 3scale integration shows up in UI
  * allow skipping apbs in the sync step
  * Merge pull request #1 from matzew/Docker_Horst
  * Adding :whale:host where needed, since on Fedora that's not defaulting to docker.io
  * pull images from feedhenry repo to target repo
  * Merge pull request #141 from feedhenry/readme-updates
  * Some readme updates to remove override for dockerhub_org
  * Remove dockerhub_org param from getting started docs. Defaults to feedhenry - can still be overriden
  * Merge pull request #140 from feedhenry/mobile-icon-padding
  * Updated mobile icon layout
  * Merge pull request #139 from sedroche/update-to-modal
  * Add service modal to get started screen
  * Merge pull request #138 from sedroche/download-url
  * Fix missed renaming
  * Set content headers on download response
  * Add build download url
  * Merge pull request #137 from feedhenry/build-history-empty-state
  * Only show Build History tab if a Build Config was created
  * Merge pull request #136 from feedhenry/update-buildfarm-name
  * Update buildfarm name to 'aerogear-digger'
  * Merge pull request #130 from feedhenry/add-insecure-option-to-template
  * Merge pull request #126 from sedroche/build-configs
  * Merge pull request #131 from matzew/add_token_to_gitignore
  * Merge pull request #134 from matzew/Travis_go_19
  * Go 1.9 in Travis
  * correcting token
  * Merge pull request #133 from matzew/Remove_Fake_artifacts
  * :boom: Removing fake templates
  * accept insecure as flag and set in template params. Add log level too
  * Review feedback
  * Rely on modal library to position modal
  * Revert "🎨 Move the external service modal out of the main content div"
  * Fix some UI bugs
  * Move components to mobile-app directory
  * Add create config component and UI fixes
  * Filter Builds by the BuildConfig name
  * Initial app build flow
  * add error for bad response code
  * Merge pull request #129 from feedhenry/walkthroughs
  * Various changes to the walkthrough
  * Merge pull request #127 from matzew/Metrics_try
  * Adding (optional) Hawkular CPU/mem metrics
  * Merge pull request #125 from feedhenry/FH-4215-temp-download
  * Download artifact from Jenkins
  * walkthrough for none local dev setup
  * Merge pull request #128 from matzew/More_retries
  * Moar retries for the MCP extension
  * Merge pull request #114 from sedroche/service-create-into-modal
  * Use N/A, same as help text
  * 🎨 Move the external service modal out of the main content div
  * allow build downloads to be created
  * Improve namespace messaging
  * Merge pull request #123 from feedhenry/FH-4133-build-resources
  * add support for a password and p12
  * Merge pull request #119 from feedhenry/FH-4133-build-resources
  * Merge pull request #122 from feedhenry/server-improvements
  * remove panics and replace with log fatal
  * add sigterm to signals
  * some minor improvements around timeouts
  * improve test
  * allow uploading credentials
  * Merge pull request #110 from rtbm/local_bower
  * Merge branch 'master' into local_bower
  * Merge pull request #111 from rtbm/local_grunt
  * Merge pull request #116 from wei-lee/fix-building-cli-error
  * Merge pull request #115 from feedhenry/FH-4189-git-source-secret
  * Merge pull request #120 from feedhenry/fix-asb-install
  * Merge pull request #121 from wei-lee/fix-service-integration-typo
  * fix asb install
  *  🐛 fix the typo that makes it impossible to create integrations for  mobile services
  * add quotes
  * Merge pull request #117 from feedhenry/fix-asb-install
  * add fixed version of image rather than latest
  *  ⬆ remove the deps that are suggested by `dep ensure`
  *  🐛 add the missing deps file and fix the `make build_cli` error
  * update failing test
  * fix json tag
  * remove invali char
  * add generatekeys endpoint
  * Merge pull request #107 from feedhenry/FH-4128-create-build
  * break out the creation of the source secret. Add more tests
  * Put service creation into a modal
  * WIP first pass at creating the required build configs
  * Locally installed grunt build
  * Locally installed bower build

0.0.2 / 2017-09-28
==================

  * Merge pull request #104 from matzew/Change_to_loopback_IP
  * Merge pull request #105 from aidenkeating/FH-4055-keycloak-integration-details
  * Merge pull request #106 from matzew/License_Header
  * Fixing incorrect license header
  * FH-4055 Add Keycloak integration details
  * localhost does not work, it does need the IP address...
  * Merge pull request #103 from feedhenry/FH-4139-partial-mock-sync-dashboard-review
  * Add reviewed keycloak dashboard UI
  * Merge pull request #102 from philbrookes/FH-4081
  * use cordova icon
  * Merge pull request #100 from feedhenry/fake-templates
  * Checkpoint for hooking up fh-sync-server reviewed dashboard to existing & new metrics being gathered
  * add fake services for showing concept demo
  * Merge pull request #99 from philbrookes/FH-4081
  * add icons to APBs
  * Merge pull request #97 from aidenkeating/FH-4070-apiKey-integration
  * FH-4070 Move all API Keys into JSON value
  * Checkpoint UI changes for sync dashboard review showing workers & queues & timings
  * Merge pull request #98 from feedhenry/add-license-1
  * Create LICENSE
  * Merge pull request #86 from sedroche/complete-mobileoverview-component
  * FH-4070 Remove duplicate error checking
  * FH-4070 Move API Key Client to Mobile App
  * FH-4070 Integrate API Keys with services
  * Merge pull request #96 from feedhenry/move-clients
  * gofmt
  * rename pkg to make it clearer what it is
  * rename pkg to make it clearer what it is
  * Merge pull request #95 from feedhenry/remove-mcp-specific-apb
  * remove import
  * Merge branch 'master' of github.com:feedhenry/mcp-standalone into remove-mcp-specific-apb
  * Merge pull request #91 from philbrookes/FH-4058
  * fix unit tests
  * fix issue with ui
  * remove mcp specific apb
  * refactor mobile.User and re-use mocks
  * refactor mobile.User and re-use mocks
  * fix unit tests
  * unit tests for authChecker
  * inject externalHTTPRequester
  * changes per review
  * remove fmt output
  * Fix unit tests
  * show message if service is not writeable
  * gofmt
  * implement auth checker
  * working on authChecker
  * update unit tests
  * deconfigure
  * cross-namespace secret mounting
  * Merge pull request #93 from feedhenry/factor-out-mobile-service-repo
  * removed token client builder
  * remove service repo builder from token client builder
  * start removing the tokenclientbuilder
  * start removing the tokenclientbuilder
  * Merge pull request #89 from aidenkeating/FH-4070-remove-api-keys-on-delete
  * FH-4070 Match comment with function name
  * Merge pull request #84 from sedroche/add-cert-troubleshooting
  * FH-4070 Remove API Keys on app delete
  * Merge pull request #88 from feedhenry/move-oc-data-dirs
  * 🐳 Use paths under ./ui for OpenShift data dirs
  * Merge pull request #87 from feedhenry/ignore-ide-folders
  * ignore Gogland & VSCode
  * Merge pull request #80 from aidenkeating/FH-4070-api-keys-configmap
  * FH-4070 Move map creation to Mobile App Repo
  * FH-4070 Make default error code to 500
  * FH-4070 Ensure secret exists on server start
  * FH-4070 Move app types to constants
  * FH-4070 Update status code on invalid clientType
  * FH-4070-create-api-keys-configmap
  * FH-4070 Add mcp-mobile-keys configmap on MCP create
  * Add mobile overview component
  * Merge pull request #82 from feedhenry/server-docs
  * updates after review
  * Add Linux cert cache clearing
  * Merge pull request #75 from matzew/SelfSigned_Certs
  * Merge pull request #81 from sedroche/add-overview-card
  * Add overview component
  * Merge pull request #85 from philbrookes/change-installer-asb-file
  * change path to asb template file
  * Merge pull request #78 from feedhenry/FH-4027-metrics
  * Split up sync metrics struct into smaller re-usable parts. Use a regex for parsing string metrics into int64. Some tidy up based on feedback
  * first pass at a doc around code and design
  * Add unit test for sync server metrics gatherer
  * A little bit more on making the certificate issue more visible, also added notes on a potential MCP provisioning issue
  * Merge pull request #77 from sedroche/fix-icon
  * Merge pull request #79 from philbrookes/mobileservice-duplicate-data
  * add in panic recover for background process
  * remove duplicate data
  * Merge branch 'FH-4027-metrics' of github.com:feedhenry/mcp-standalone into FH-4027-metrics
  * merge in master. add keycloack metrics and metrics background gatherer
  * add keycloack metrics. metrics runner and tests
  * Update mobile icon classes
  * Merge pull request #72 from feedhenry/service-delete
  * add delete mobile service
  * Add empty message if no stats are available for a service
  * Add sync metrics debug logging
  * Add fh-sync-server stats gatherer
  * Merge pull request #76 from sedroche/update-cards
  * Add delete functionality to app/services
  * incremental
  * Merge branch 'FH-4027-metrics' of github.com:feedhenry/mcp-standalone into FH-4027-metrics
  * Merge pull request #74 from feedhenry/remove-oauth-handler
  * Merge pull request #73 from philbrookes/FH-4082
  * Add circle to service cards
  * Merge pull request #46 from aidenkeating/update-getting-started-sync-keycloak
  * remove the oauth handler as it is no longer needed
  * Update examples in Sync/Keycloak getting started
  * Show icons in mobile overview
  * incremental
  * Merge pull request #71 from matzew/Improve_doc
  * :lipstick: some little getting started improvements
  * Chart Title style fixup
  * Merge pull request #54 from feedhenry/mcp-onboarding
  * Merge pull request #69 from aidenkeating/fix-cordova-config-template
  * Merge pull request #70 from philbrookes/FH-4082
  * add mobile tags for catalog categories
  * Wait for postdigest before rendering charts
  * Add Chart rendering on Service Dashboard
  * add initial metrics endpoint
  * add initial metrics endpoint
  * add initial metrics endpoint
  * Merge pull request #68 from philbrookes/FH-4068
  * remove MCP_URL caching
  * Merge pull request #66 from sedroche/tabs-view
  * Merge pull request #67 from sedroche/highlight-tab-for-service
  * initial loop setup
  * Fix cordova config template
  * Include service prefix to keep tab highlighting
  * Only show tabs for > 1 app
  * Merge pull request #64 from sedroche/onboarding-fixes
  * Merge pull request #65 from philbrookes/FH-4074
  * split integrations into 2 tabs
  * Point to correct guides, networking fix and troubleshooting
  * Merge pull request #63 from feedhenry/FH-4058-external-service-ns
  * add ns to local services
  * add namespace dynamically to non external services
  * add namespace to service type
  * Merge pull request #62 from feedhenry/refactor-pkgs
  * changes after review
  * gofmt
  * Merge pull request #60 from philbrookes/FH-4053
  * refactoring
  * update installation field name
  * Merge pull request #58 from philbrookes/FH-4053
  * update keycloak secret name/type and fix some minor code issues
  * Merge pull request #59 from sedroche/onboarding-fixes
  * README changes from onboarding
  * Do bower install as part of install script
  * Update lock file after insalling dependencies
  * Merge pull request #57 from philbrookes/FH-4041-remove-configuration
  * add missing validation
  * Merge pull request #56 from feedhenry/delete-web-folder-again
  * remove remnants of web dir
  * Merge pull request #55 from feedhenry/FH-4033
  * update the apbs to use the authorization header. Fix some minor naming issues
  * update ansible service broker and fix installation to use admin role
  * Merge pull request #53 from philbrookes/FH-4041-remove-configuration
  * validate incoming params are not empty
  * update missing docs
  * update tests
  * change mocked field name from name to type
  * 📝 Add extra Onboarding docs with links to various resources
  * reduce duplication in test
  * update docs
  * setup dist folder in ui ansible-playbook
  * update ui with new end-point
  * adding configure / deconfigure endpoints
  * Merge pull request #52 from feedhenry/external-services
  * Merge branch 'master' into external-services
  * Merge pull request #51 from feedhenry/FH-3812-show-url-and-login-credentials-for-services
  * fix failing tests
  * ✨ FH-3812 Show url & credentials on service dashboard
  * adds the capability to add external services.
  * Merge branch 'master' of github.com:feedhenry/mcp-standalone into external-services
  * Merge pull request #33 from feedhenry/basic-cli
  * remove cli build
  * bump go version
  * update readme
  * update gopkg
  * update to master
  * resolve conflict
  * incremental
  * Merge pull request #49 from feedhenry/empty-states-on-mobile-overview-page
  * 💄 Add empty state UI for the various states on the Mobile Overview page
  * resolve conflicts
  * incremental
  * Merge pull request #48 from feedhenry/fixup-mobile-apps-page
  * Various style adjustments and minor fixes for views...
  * resolve conflict
  * add a mobile service create endpoint
  * Merge pull request #47 from feedhenry/FH-4010-add-sync-service-dashboard
  * Fixup styles around overview page & service screen
  * Merge pull request #45 from feedhenry/remove-deprecated-web
  * Remove deprecated web folder in favour of ui extension in ui folder
  * increment
  * Merge branch 'master' of github.com:feedhenry/mcp-standalone into external-services
  * incremental
  * merge in master
  * add cli deps
  * build cli as part of build
  * very basic cli for dealing with mobileapps via the api

before-ui-extension / 2017-09-07
================================

  * Merge pull request #44 from feedhenry/add-integration-screen
  * ✨ Add the integrations screen
  * Merge pull request #43 from feedhenry/readme
  * add readme change for create app
  * Merge pull request #42 from feedhenry/FH-3826-bundle-js-and-css-via-grunt
  * Additional local dev fixups...
  * Merge pull request #40 from feedhenry/philbrookes-add-mobileservices-routes
  * remote prints
  * Merge branch 'master' into FH-3826-bundle-js-and-css-via-grunt
  * FH-3826 Bundle js & css for the mcp ui extension
  * ✨ FH-4023 Import the UI extension code & installation setup
  * fix config test
  * add integration view for sync sync with keycloak
  * Merge pull request #35 from feedhenry/philbrookes-add-mobileservices-routes
  * Merge pull request #39 from philbrookes/add-walkthrough
  * Merge branch 'master' of github.com:feedhenry/mcp-standalone into philbrookes-add-mobileservices-routes
  * incremental
  * Merge pull request #38 from feedhenry/FH-4023-add-mcp-ui-extension
  * add walkthrough
  * ✨ FH-4023 Import the UI extension code & installation setup
  * Merge pull request #37 from philbrookes/add-installer
  * add installer
  * Merge pull request #36 from aidenkeating/allow-x-app-api-key-header
  * Include x-app-api-key in allowed headers
  * merge in master
  * incremental
  * Merge pull request #32 from feedhenry/FH-3972-mobile-app-view
  * 🎨 Integration config steps fixups
  * Merge branch 'FH-3972-mobile-app-view' of github.com:feedhenry/mcp-standalone into FH-3972-mobile-app-view
  * merge in master
  * updates to service handler
  * 🔧 FH-4008 Get OAuth Scope from the server
  * fix bad merge
  * merge in master changes
  * Merge pull request #34 from aidenkeating/include-more-keycloak-config-values
  * Include clientId and url values in KeycloakConfig
  * pass at the mobile app view
  * add mobile apps view and hook up to api. Also show mobile service. Add start of app view and integration view
  * Merge pull request #30 from feedhenry/FH-4002-set-oauth-config-dynamically-based-on-request
  * 🔧 FH-4002 Set the oauth redirect config dynamically based on the Request headers
  * Merge pull request #26 from feedhenry/FH-3876
  * add missing headers helper
  * updates after code review
  * add mobileservice configure end-point and unit tests
  * Merge pull request #27 from feedhenry/FH-3990-call-to-openshift-on-startup
  * 🔧 FH-3990 Call to OpenShift on startup to get server metadata
  * Update README.md
  * remove unused code
  * initial pass at changes to allow for different service configuration
  * Merge pull request #25 from feedhenry/change-fh-sync-to-fh-sync-server
  * change fh-sync to fh-sync-server
  * Merge pull request #22 from feedhenry/update-origin-web-common-with-scope-change
  * Merge pull request #23 from feedhenry/apbs
  * fix logger error
  * add apbs for iOS cordova mcp and android
  * 🐛 Use upstream version of origin-web-common with scope feature
  * Merge pull request #21 from feedhenry/FH-3970-use-patternfly
  * 💄 FH-3970 Use Patternfly
  * 🎨 Use `grunt serve` as the UI server for local dev
  * Merge pull request #17 from feedhenry/create-role-binding
  * changes after code review
  * make using sa token secure. Add skip rolebinding optional header
  * update the config and makefile to use namespace mcp-standalone by default
  * update to makefile and readme
  * add a make run command with sane defaults
  * fix interface type
  * add rolebinding mw to work around apb limititation
  * Merge pull request #18 from feedhenry/add-mobileservices-routes
  * add mobileservices routes
  * Merge pull request #16 from feedhenry/temp-fix-for-upstream-dep
  * 🐛 Temporary fix to include fork of origin-web-common with Scope fix
  * Merge pull request #15 from feedhenry/use-web-app-src-as-default
  * 🐛 Use web/app src for local & web/dist for build/docker
  * Merge pull request #14 from feedhenry/refactor-new-repo
  * fix up imports and names for new repo
  * changes to support https serving
  * Merge pull request #2 from feedhenry/add-initial-web-ui
  * ✨ Add Initial AngularJS based Web UI
  * add make and install
  * Merge pull request #8 from feedhenry/build-improvements
  * add new make targets for gofmt and lint
  * Merge pull request #7 from feedhenry/metrics
  * add metrics to the server
  * Merge pull request #6 from feedhenry/mobile-services
  * add a mobile service handler to list mobile services
  * Merge pull request #5 from feedhenry/update-to-template
  * add readiness and liveness probes. Wire up sys route
  * Merge pull request #4 from feedhenry/more-unit-tests
  * add sdk config handler test and mobile service repo test
  * Merge pull request #3 from feedhenry/add-sdk-config
  * refactor and add in sdk handler
  * Merge pull request #1 from feedhenry/test-travis
  * add basic readme
  * Add tests and various updates. Add travis
  * vendor
  * incremental
