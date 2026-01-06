# PHP Examples

Complete examples of PHP project analysis and deployment configuration.

## Example 1: Laravel Application

**User request**: "Pack this Laravel project"

**Detected files**:
- `composer.json` with `laravel/framework`
- `artisan` (Laravel CLI)
- `composer.lock`
- `.env.example`

**Output**:
```json
[
  {
    "language": "php",
    "version": "8.2",
    "apt": [],
    "dev": {
      "environment": {
        "APP_ENV": "local",
        "APP_DEBUG": "true"
      },
      "setup": ["composer install"],
      "commands": ["php artisan serve --host=0.0.0.0 --port=8080"]
    },
    "prod": {
      "environment": {
        "APP_ENV": "production",
        "APP_DEBUG": "false"
      },
      "setup": [
        "composer install --no-dev --optimize-autoloader",
        "php artisan config:cache",
        "php artisan route:cache",
        "php artisan view:cache"
      ],
      "commands": ["php artisan serve --host=0.0.0.0 --port=8080"]
    },
    "port": 8080,
    "entrypoint": "#!/bin/bash\n\napp_env=${1:-development}\n\nif [ \"$app_env\" = \"production\" ] || [ \"$app_env\" = \"prod\" ] ; then\n    export APP_ENV=production\n    export APP_DEBUG=false\n    composer install --no-dev --optimize-autoloader\n    php artisan config:cache\n    php artisan route:cache\n    php artisan view:cache\n    php artisan serve --host=0.0.0.0 --port=8080\nelse\n    export APP_ENV=local\n    export APP_DEBUG=true\n    composer install\n    php artisan serve --host=0.0.0.0 --port=8080\nfi",
    "evidence": {
      "files": ["composer.json", "artisan", "composer.lock"],
      "reason": "PHP Laravel application detected. Using artisan serve (production typically uses nginx/php-fpm)."
    }
  }
]
```

**Key Points**:
- PHP has **highest priority (#1)** in language detection
- Laravel detected from `laravel/framework` in composer.json
- Dev environment enables debug mode
- Production optimizes autoloader and caches config/routes/views
- `artisan serve` is development server (production usually needs nginx/php-fpm)
- No APT packages needed - PHP base images include extensions

---

## Example 2: Symfony Application

**User request**: "Analyze Symfony project"

**Detected files**:
- `composer.json` with `symfony/framework-bundle`
- `symfony.lock`
- `bin/console`

**Output**:
```json
[
  {
    "language": "php",
    "version": "8.2",
    "apt": [],
    "dev": {
      "environment": {
        "APP_ENV": "dev",
        "APP_DEBUG": "1"
      },
      "setup": ["composer install"],
      "commands": ["php -S 0.0.0.0:8080 -t public"]
    },
    "prod": {
      "environment": {
        "APP_ENV": "prod",
        "APP_DEBUG": "0"
      },
      "setup": [
        "composer install --no-dev --optimize-autoloader",
        "php bin/console cache:clear",
        "php bin/console cache:warmup"
      ],
      "commands": ["php -S 0.0.0.0:8080 -t public"]
    },
    "port": 8080,
    "entrypoint": "#!/bin/bash\n\napp_env=${1:-development}\n\nif [ \"$app_env\" = \"production\" ] || [ \"$app_env\" = \"prod\" ] ; then\n    export APP_ENV=prod\n    export APP_DEBUG=0\n    composer install --no-dev --optimize-autoloader\n    php bin/console cache:clear\n    php bin/console cache:warmup\n    php -S 0.0.0.0:8080 -t public\nelse\n    export APP_ENV=dev\n    export APP_DEBUG=1\n    composer install\n    php -S 0.0.0.0:8080 -t public\nfi",
    "evidence": {
      "files": ["composer.json", "bin/console", "symfony.lock"],
      "reason": "Symfony framework detected. Using PHP built-in server for simplicity."
    }
  }
]
```

**Key Points**:
- Symfony detected from `symfony/framework-bundle`
- Console commands for cache management
- Built-in PHP server: `php -S 0.0.0.0:8080 -t public`
- Production disables debug and optimizes cache

---

## Example 3: Plain PHP Application

**User request**: "Pack PHP app"

**Detected files**:
- `index.php`
- `composer.json` (minimal, no framework)
- PHP files in root

**Output**:
```json
[
  {
    "language": "php",
    "version": "8.1",
    "apt": [],
    "dev": {
      "environment": {},
      "setup": ["composer install"],
      "commands": ["php -S 0.0.0.0:8080"]
    },
    "prod": {
      "environment": {},
      "setup": ["composer install --no-dev"],
      "commands": ["php -S 0.0.0.0:8080"]
    },
    "port": 8080,
    "entrypoint": "#!/bin/bash\n\napp_env=${1:-development}\n\nif [ \"$app_env\" = \"production\" ] || [ \"$app_env\" = \"prod\" ] ; then\n    composer install --no-dev\n    php -S 0.0.0.0:8080\nelse\n    composer install\n    php -S 0.0.0.0:8080\nfi",
    "evidence": {
      "files": ["index.php", "composer.json"],
      "reason": "Plain PHP application detected. Using built-in PHP development server."
    }
  }
]
```

---

## PHP Detection Patterns

### Required Files
- `composer.json` (REQUIRED for dependency management)
- `*.php` files
- `composer.lock` (optional, dependency lock)

### Framework Detection
| Framework | Detection | Default Port |
|-----------|-----------|--------------|
| Laravel | `laravel/framework` in composer.json, `artisan` file | 8080 |
| Symfony | `symfony/framework-bundle` in composer.json, `bin/console` | 8080 |
| CodeIgniter | `codeigniter4/framework` in composer.json | 8080 |
| Yii | `yiisoft/yii2` in composer.json | 8080 |
| Slim | `slim/slim` in composer.json | 8080 |
| Plain PHP | composer.json without framework deps | 8080 |

### Host Binding Syntax
- Laravel artisan: `--host=0.0.0.0 --port=PORT`
- PHP built-in server: `-S 0.0.0.0:PORT`
- Symfony: `-S 0.0.0.0:PORT -t public`

### Version Detection
1. `composer.json` → `config.platform.php` field
2. `.php-version` file (if exists)
3. Omit if uncertain

### Composer Strategies

**Development**:
```bash
composer install
```

**Production**:
```bash
composer install --no-dev --optimize-autoloader
```

### Common Environment Variables
- Laravel:
  - `APP_ENV`: `local` | `production`
  - `APP_DEBUG`: `true` | `false`
  - `APP_KEY`: Application key (required)
- Symfony:
  - `APP_ENV`: `dev` | `prod`
  - `APP_DEBUG`: `1` | `0`
  - `APP_SECRET`: Application secret

### Laravel-Specific Commands

**Development**:
- `php artisan serve`: Start dev server
- `php artisan migrate`: Run database migrations

**Production**:
- `php artisan config:cache`: Cache configuration
- `php artisan route:cache`: Cache routes
- `php artisan view:cache`: Cache views
- `php artisan optimize`: Optimize framework

### Symfony-Specific Commands

**Development**:
- `php bin/console server:start`: Start dev server (older versions)
- `symfony server:start`: Symfony CLI server

**Production**:
- `php bin/console cache:clear`: Clear cache
- `php bin/console cache:warmup`: Warm up cache
- `php bin/console assets:install`: Install assets

### APT Dependencies
**Default: `[]` (empty array)**

Reasons:
- PHP base images include common extensions
- Most PHP extensions available as Docker layers
- Composer handles PHP-level dependencies

**Only include APT packages if**:
- Project explicitly requires system libraries
- Database drivers not included in base image
- Image processing tools: `libpng-dev`, `libjpeg-dev`
- CLI tools: `postgresql-client`

### Common Extensions (Usually Pre-installed)
- `pdo_mysql`, `pdo_pgsql`: Database drivers
- `gd`: Image processing
- `curl`: HTTP client
- `mbstring`: Multibyte string support
- `xml`: XML processing
- `zip`: Archive handling
- `redis`: Redis support

### Default Port
**8080** (standard for PHP applications)

### Production Considerations
1. **Web Server**: Built-in PHP server is for development
   - Production should use nginx + php-fpm or Apache + mod_php
   - Artisan serve is not recommended for production
2. **Optimization**: Use `--optimize-autoloader` for composer
3. **Caching**: Cache configuration, routes, and views
4. **Environment**: Set `APP_ENV=production` and disable debug
5. **Assets**: Compile and minify frontend assets
