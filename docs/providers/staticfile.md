# Staticfile Provider

- Detection: Uses confidence-based detection with weighted indicators and exclusion logic:
  - HTML files: `*.html`, `*.htm` (weight: 30)
  - CSS files: `*.css` (weight: 20)
  - JavaScript files: `*.js` (weight: 15)
  - Image assets: `*.png`, `*.jpg`, `*.jpeg`, `*.gif`, `*.svg` (weight: 10)
  - Entry point: `index.html` (weight: 25)
  - Configuration: `Staticfile` (weight: 15)
  - Asset directories: `assets/`, `static/`, `public/` (weight: 5)
  - Web assets: `favicon.ico`, `robots.txt` (weight: 5)
  - Detection threshold: 0.2 (20% confidence)
  - **Exclusion Logic**: Confidence reduced by 70% if other language files detected (package.json, composer.json, requirements.txt, go.mod, Cargo.toml, Gemfile, pom.xml, build.gradle, source code files)
  - **Pure Static Exception**: Full confidence maintained for pure HTML/CSS/JS projects without other language indicators
  - **Special Cases**: Full confidence for projects with only `index.html` + `test.json`, or projects with `Staticfile` configuration

- Version Detection: Static file projects use version `latest` by default
  - No specific version detection as static files are environment-independent

- Framework Detection: No framework detection (static files are framework-agnostic)

- Package Manager: No package manager (static files are self-contained)

- Commands:
  - **Development**: 
    - No development command (static files don't require compilation)
  - **Build**: 
    - No build command (static files are served as-is)
  - **Start**: 
    - `nginx -g 'daemon off;'`

- Native Compilation Detection: Static files don't require compilation
  - Returns `false` as static files are served directly

- Metadata: Provides comprehensive metadata including:
  - `hasHtml`: Presence of HTML files
  - `hasCss`: Presence of CSS files
  - `hasJs`: Presence of JavaScript files
  - `hasImages`: Presence of image assets
  - `hasIndexHtml`: Presence of `index.html` entry point
  - `hasStaticfile`: Presence of `Staticfile` configuration
  - `hasAssets`: Presence of asset directories