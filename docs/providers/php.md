# PHP Provider

- Detection: Uses confidence-based detection with weighted indicators:
  - Composer configuration: `composer.json` (weight: 30)
  - PHP source files: `*.php` (weight: 25)
  - Composer lock file: `composer.lock` (weight: 15)
  - Entry point files: `index.php`, `app.php` (weight: 10)
  - Dependencies: `vendor/`, `autoload.php` (weight: 10)
  - Version/config files: `.php-version`, `phpunit.xml` (weight: 5)
  - Framework files: `artisan`, `wp-config.php` (weight: 5)
  - Detection threshold: 0.2 (20% confidence)

- Version Detection: Priority order for PHP version resolution:
  1. `composer.json` require.php field
  2. `.php-version` file
  3. Default version: `8.2`

- Framework Detection: Automatically detects popular PHP frameworks by analyzing Composer dependencies:
  - **MVC Frameworks**: Laravel, Symfony, CodeIgniter, CakePHP, Yii2
  - **Micro Frameworks**: Slim Framework
  - **Legacy**: Zend Framework, Laminas, Phalcon
  - **ORM/Template**: Doctrine ORM, Twig
  - **Analysis**: Checks `composer.json` require dependencies

- Package Manager: Always uses Composer for dependency management

- Commands:
  - **Development**: 
    - `composer install`
    - **With index.php**: `php -S 0.0.0.0:8000 index.php`
    - **Without index.php**: `php -S 0.0.0.0:8000`
  - **Build**: 
    - `composer install --no-dev --optimize-autoloader`
  - **Start**: 
    - **With index.php**: `php -S 0.0.0.0:8000 index.php`
    - **Without index.php**: `php -S 0.0.0.0:8000`

- Native Compilation Detection: PHP projects typically don't require native compilation
  - Returns `false` as most PHP projects are interpreted

- Metadata: Provides comprehensive metadata including:
  - `hasComposerJson`: Presence of `composer.json`
  - `hasComposerLock`: Presence of `composer.lock`
  - `hasPHPSrc`: Presence of PHP source files
  - `hasIndex`: Presence of `index.php`
  - `hasVendor`: Presence of `vendor` directory
  - `framework`: Detected framework name