# Python Provider

- Detection: Uses confidence-based detection with weighted indicators:
  - Dependency files: `requirements.txt`, `pyproject.toml`, `setup.py`, `Pipfile` (weight: 30)
  - Python source files: `*.py` (weight: 25)
  - Lock files: `poetry.lock` (weight: 15), `Pipfile.lock` (weight: 15)
  - Version files: `.python-version`, `runtime.txt` (weight: 10)
  - Entry point files: `manage.py`, `app.py`, `main.py` (weight: 5)
  - Python artifacts: `__pycache__`, `*.pyc` (weight: 5)
  - Detection threshold: 0.3 (30% confidence)

- Version Detection: Priority order for Python version resolution:
  1. `.python-version` file
  2. `runtime.txt` file (Heroku format: `python-X.Y.Z`)
  3. `pyproject.toml` python requirement
  4. `Pipfile` python_version field
  5. Default version: `3.11`

- Framework Detection: Automatically detects popular Python frameworks by analyzing dependency files:
  - **Web Frameworks**: Django, Flask, FastAPI, Tornado, Pyramid, Bottle, Sanic, Quart, Starlette
  - **Data Science**: Streamlit, Dash, Jupyter
  - **Analysis**: Checks `requirements.txt`, `pyproject.toml`, and `Pipfile` for framework dependencies

- Package Manager Detection: Automatically detects package manager based on project files:
  - `poetry.lock` → poetry
  - `Pipfile.lock` → pipenv
  - `requirements.txt` → pip
  - `pyproject.toml` → pip (default)

- Commands:
  - **Development**: 
    - **requirements.txt**: `pip install -r requirements.txt`, then start command
    - **Pipfile**: `pipenv install --dev`, then start command
    - **pyproject.toml**: `pip install -e .`, then start command
  - **Build**: 
    - **requirements.txt**: `pip install -r requirements.txt`
    - **Pipfile**: `pipenv install`
    - **pyproject.toml**: `pip install .`
  - **Start**: 
    - **app.py**: `python app.py`
    - **main.py**: `python main.py`
    - **manage.py** (Django): `python manage.py runserver 0.0.0.0:8000`
    - **Fallback**: `python -m http.server 8000`

- Native Compilation Detection: Automatically detects packages requiring native compilation:
  - **Common packages**: numpy, scipy, pandas, pillow, lxml, psycopg2, mysqlclient, cryptography, cffi, cython, pycrypto, gevent

- Metadata: Provides comprehensive metadata including:
  - `hasRequirements`: Presence of `requirements.txt`
  - `hasPyprojectToml`: Presence of `pyproject.toml`
  - `hasSetupPy`: Presence of `setup.py`
  - `hasPipfile`: Presence of `Pipfile`
  - `hasPoetryLock`: Presence of `poetry.lock`
  - `packageManager`: Detected package manager
  - `framework`: Detected framework name