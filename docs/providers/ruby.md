# Ruby Provider

- Detection: Uses confidence-based detection with weighted indicators:
  - Gemfile: `Gemfile` (weight: 30)
  - Lock file: `Gemfile.lock` (weight: 20)
  - Ruby source files: `*.rb` (weight: 20)
  - Version files: `.ruby-version`, `.rvmrc` (weight: 10)
  - Configuration files: `config.ru`, `Rakefile` (weight: 10)
  - Directory structure: `app/`, `lib/` (weight: 5)
  - Test directories: `spec/`, `test/` (weight: 5)
  - Detection threshold: 0.3 (30% confidence)

- Version Detection: Priority order for Ruby version resolution:
  1. `.ruby-version` file
  2. `.rvmrc` file (ruby-X.Y.Z format)
  3. `Gemfile` ruby directive
  4. Default version: `3.2`

- Framework Detection: Automatically detects popular Ruby frameworks by analyzing Gemfile:
  - **Web Frameworks**: Ruby on Rails, Sinatra, Grape, Hanami, Roda, Cuba, Padrino
  - **Static Site Generators**: Jekyll, Middleman
  - **Analysis**: Checks `Gemfile` for framework gem declarations

- Package Manager: Always uses Bundler for dependency management

- Commands:
  - **Development**: 
    - `bundle install`
    - **With config.ru**: `bundle exec rackup -o 0.0.0.0 -p 4567`
    - **With app.rb**: `ruby app.rb`
    - **Fallback**: `ruby -run -e httpd . -p 4567`
  - **Build**: 
    - `bundle install --without development test`
  - **Start**: 
    - **With config.ru**: `bundle exec rackup -o 0.0.0.0 -p 4567`
    - **With app.rb**: `ruby app.rb`
    - **Fallback**: `ruby -run -e httpd . -p 4567`

- Native Compilation Detection: Ruby projects typically don't require native compilation
  - Returns `false` unless gems have C extensions (not currently detected)

- Metadata: Provides comprehensive metadata including:
  - `hasGemfile`: Presence of `Gemfile`
  - `hasGemfileLock`: Presence of `Gemfile.lock`
  - `hasRakefile`: Presence of `Rakefile`
  - `hasConfigRu`: Presence of `config.ru`
  - `hasRubyVersion`: Presence of `.ruby-version`
  - `framework`: Detected framework name