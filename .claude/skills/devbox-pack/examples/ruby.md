# Ruby Examples

Complete examples of Ruby project analysis and deployment configuration.

## Example 1: Ruby on Rails Application

**User request**: "Analyze Rails app"

**Detected files**:
- `Gemfile` with `rails 7.0`
- `config/application.rb`
- `.ruby-version` → 3.2

**Output**:
```json
[
  {
    "language": "ruby",
    "version": "3.2",
    "apt": [],
    "dev": {
      "environment": {},
      "setup": ["bundle install"],
      "commands": ["bundle exec rails server -b 0.0.0.0 -p 3000"]
    },
    "prod": {
      "environment": {
        "RAILS_ENV": "production"
      },
      "setup": [
        "bundle install --without development test",
        "bundle exec rails db:migrate"
      ],
      "commands": ["bundle exec rails server -b 0.0.0.0 -p 3000"]
    },
    "port": 3000,
    "entrypoint": "#!/bin/bash\n\napp_env=${1:-development}\n\nif [ \"$app_env\" = \"production\" ] || [ \"$app_env\" = \"prod\" ] ; then\n    export RAILS_ENV=production\n    bundle install --without development test\n    bundle exec rails db:migrate\n    bundle exec rails server -b 0.0.0.0 -p 3000\nelse\n    bundle install\n    bundle exec rails server -b 0.0.0.0 -p 3000\nfi",
    "evidence": {
      "files": ["Gemfile", "config/application.rb", ".ruby-version"],
      "reason": "Ruby on Rails 7.0 application detected"
    }
  }
]
```

**Key Points**:
- Rails-specific commands: `db:migrate`
- Bundle install with groups: `--without development test` in prod
- Version from `.ruby-version` file
- `RAILS_ENV` environment variable
- Default Rails port: 3000

---

## Example 2: Sinatra API Application

**User request**: "Pack Sinatra API"

**Detected files**:
- `Gemfile` with `sinatra`
- `app.rb` or `config.ru`
- `Gemfile.lock`

**Output**:
```json
[
  {
    "language": "ruby",
    "version": "3.1",
    "apt": [],
    "dev": {
      "environment": {
        "RACK_ENV": "development"
      },
      "setup": ["bundle install"],
      "commands": ["bundle exec ruby app.rb -o 0.0.0.0 -p 4567"]
    },
    "prod": {
      "environment": {
        "RACK_ENV": "production"
      },
      "setup": ["bundle install --without development test"],
      "commands": ["bundle exec rackup -o 0.0.0.0 -p 4567"]
    },
    "port": 4567,
    "entrypoint": "#!/bin/bash\n\napp_env=${1:-development}\n\nif [ \"$app_env\" = \"production\" ] || [ \"$app_env\" = \"prod\" ] ; then\n    export RACK_ENV=production\n    bundle install --without development test\n    bundle exec rackup -o 0.0.0.0 -p 4567\nelse\n    export RACK_ENV=development\n    bundle install\n    bundle exec ruby app.rb -o 0.0.0.0 -p 4567\nfi",
    "evidence": {
      "files": ["Gemfile", "app.rb", "Gemfile.lock"],
      "reason": "Sinatra application detected. Default Sinatra port is 4567."
    }
  }
]
```

**Key Points**:
- Sinatra default port: 4567
- Environment: `RACK_ENV` instead of `RAILS_ENV`
- Dev: Direct Ruby execution
- Prod: Rackup for Rack applications
- Host binding: `-o 0.0.0.0`

---

## Example 3: Rails API with Puma Server

**User request**: "Analyze Rails API with Puma"

**Detected files**:
- `Gemfile` with `rails` and `puma`
- `config/puma.rb`
- `config/application.rb`

**Output**:
```json
[
  {
    "language": "ruby",
    "version": "3.2",
    "apt": [],
    "dev": {
      "environment": {
        "RAILS_ENV": "development"
      },
      "setup": ["bundle install"],
      "commands": ["bundle exec puma -C config/puma.rb -b tcp://0.0.0.0:3000"]
    },
    "prod": {
      "environment": {
        "RAILS_ENV": "production",
        "RAILS_SERVE_STATIC_FILES": "true"
      },
      "setup": [
        "bundle install --without development test",
        "bundle exec rails db:migrate",
        "bundle exec rails assets:precompile"
      ],
      "commands": ["bundle exec puma -C config/puma.rb -b tcp://0.0.0.0:3000"]
    },
    "port": 3000,
    "entrypoint": "#!/bin/bash\n\napp_env=${1:-development}\n\nif [ \"$app_env\" = \"production\" ] || [ \"$app_env\" = \"prod\" ] ; then\n    export RAILS_ENV=production\n    export RAILS_SERVE_STATIC_FILES=true\n    bundle install --without development test\n    bundle exec rails db:migrate\n    bundle exec rails assets:precompile\n    bundle exec puma -C config/puma.rb -b tcp://0.0.0.0:3000\nelse\n    export RAILS_ENV=development\n    bundle install\n    bundle exec puma -C config/puma.rb -b tcp://0.0.0.0:3000\nfi",
    "evidence": {
      "files": ["Gemfile", "config/puma.rb", "config/application.rb"],
      "reason": "Rails API with Puma server detected. Using Puma configuration file."
    }
  }
]
```

**Key Points**:
- Puma server configuration: `-C config/puma.rb`
- Asset precompilation: `rails assets:precompile` in prod
- Static file serving: `RAILS_SERVE_STATIC_FILES=true`
- Puma binding: `-b tcp://0.0.0.0:3000`

---

## Ruby Detection Patterns

### Required Files
- `Gemfile` (REQUIRED for Ruby detection)
- `Gemfile.lock` (optional, dependency lock)
- `*.rb` files

### Framework Detection
| Framework | Detection | Default Port |
|-----------|-----------|--------------|
| Rails | `rails` in Gemfile, `config/application.rb` | 3000 |
| Sinatra | `sinatra` in Gemfile, `config.ru` or `app.rb` | 4567 |
| Hanami | `hanami` in Gemfile | 2300 |
| Padrino | `padrino` in Gemfile | 3000 |
| Rack | `config.ru` only | 9292 |

### Host Binding Syntax
- Rails server: `-b 0.0.0.0`
- Puma: `-b tcp://0.0.0.0:PORT`
- Sinatra: `-o 0.0.0.0`
- Rackup: `-o 0.0.0.0`

### Version Detection
1. `.ruby-version` file (rbenv/rvm)
2. `Gemfile` → `ruby` directive
   ```ruby
   ruby '3.2.0'
   ```
3. `.tool-versions` file (asdf)
4. Omit if uncertain

### Bundle Install Strategies

**Development**:
```bash
bundle install
```

**Production**:
```bash
bundle install --without development test
# or
bundle install --deployment  # Uses Gemfile.lock versions
```

### Common Environment Variables
- `RAILS_ENV`: `development` | `production` | `test`
- `RACK_ENV`: `development` | `production`
- `RAILS_SERVE_STATIC_FILES`: `true` | `false`
- `SECRET_KEY_BASE`: Rails secret key (required in production)

### Rails-Specific Commands

**Development**:
- `bundle exec rails server`: Start dev server
- `bundle exec rails console`: Interactive console

**Production**:
- `bundle exec rails db:migrate`: Run database migrations
- `bundle exec rails assets:precompile`: Compile assets
- `bundle exec rails db:seed`: Load seed data (if needed)

### Application Servers
| Server | Usage | Notes |
|--------|-------|-------|
| WEBrick | Built-in Rails dev server | Development only |
| Puma | Modern threaded server | Rails default since 5.0 |
| Unicorn | Forking server | Legacy, still used |
| Passenger | Application server | Enterprise deployments |

### APT Dependencies
**Default: `[]` (empty array)**

Reasons:
- Ruby base images include common native extensions
- Most gems with C extensions provide prebuilt binaries
- `nokogiri`: Prebuilt binaries for common platforms
- `pg` (PostgreSQL): Binary gems available

**Include APT packages only if**:
- Project explicitly requires system libraries
- Using source compilation of native extensions
- Database client utilities needed: `postgresql-client`

### Database Gems
| Gem | Database | Binary Available? |
|-----|----------|-------------------|
| `pg` | PostgreSQL | Yes |
| `mysql2` | MySQL | Yes |
| `sqlite3` | SQLite | Yes |

### Common Ports
- Rails: 3000
- Sinatra: 4567
- Rack: 9292
- Hanami: 2300

### Production Considerations
1. **Assets**: Precompile with `rails assets:precompile`
2. **Database**: Run `rails db:migrate` before starting
3. **Secrets**: Set `SECRET_KEY_BASE` environment variable
4. **Static files**: Enable `RAILS_SERVE_STATIC_FILES` if not using nginx
5. **Dependencies**: Use `bundle install --without development test`
