# Shell Provider

- Detection: Uses confidence-based detection with weighted indicators:
  - Shell scripts: `*.sh`, `*.bash`, `*.zsh` (weight: 30)
  - Build files: `Makefile`, `makefile` (weight: 20)
  - Alternative shells: `*.fish`, `*.csh`, `*.ksh` (weight: 15)
  - Common scripts: `install.sh`, `setup.sh`, `build.sh` (weight: 10)
  - Script directories: `bin/`, `scripts/` (weight: 10)
  - Documentation: `README.md` (weight: 5)
  - Detection threshold: 0.3 (30% confidence)

- Version Detection: Shell projects use version `latest` by default
  - No specific version detection as shell scripts are environment-dependent

- Shell Type Detection: Automatically detects shell type based on file extensions:
  - **Bash**: `*.bash` files
  - **Zsh**: `*.zsh` files  
  - **Fish**: `*.fish` files
  - **Default**: `*.sh` files (defaults to `sh`)

- Project Type Detection: Categorizes shell projects by purpose:
  - **Installer**: `install.sh`, `setup.sh`, `installer.sh`
  - **Build Tool**: `build.sh`, `compile.sh`, `make.sh`
  - **Deployment**: `deploy.sh`, `release.sh`, `publish.sh`
  - **Utility**: `backup.sh`, `restore.sh`, `sync.sh`
  - **Service**: `start.sh`, `stop.sh`, `restart.sh`, `service.sh`
  - **Testing**: `test.sh`, `check.sh`, `validate.sh`
  - **Default**: `script` (generic shell scripts)

- Package Manager: No package manager (shell scripts are self-contained)

- Commands:
  - **Development**: 
    - `echo 'Shell project detected'`
  - **Build**: 
    - `echo 'No build step required for shell scripts'`
  - **Start**: 
    - `echo 'Please specify which script to run'`

- Native Compilation Detection: Shell scripts don't require compilation
  - Returns `false` as shell scripts are interpreted

- Metadata: Provides comprehensive metadata including:
  - `shellType`: Detected shell type (bash, zsh, fish, sh)
  - `projectType`: Detected project category
  - `shellFiles`: List of all shell script files
  - `hasShebang`: Whether shell files contain shebang lines