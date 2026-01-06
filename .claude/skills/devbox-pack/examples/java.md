# Java Examples

Complete examples of Java project analysis and deployment configuration.

## Example 1: Spring Boot with Maven

**User request**: "Analyze Spring Boot project"

**Detected files**:
- `pom.xml` with `spring-boot-starter-web`
- `src/main/java/`
- `mvnw` (Maven wrapper)
- `application.properties` or `application.yml`

**Output**:
```json
[
  {
    "language": "java",
    "version": "17",
    "apt": [],
    "dev": {
      "environment": {
        "SPRING_PROFILES_ACTIVE": "dev",
        "SERVER_ADDRESS": "0.0.0.0"
      },
      "commands": ["./mvnw spring-boot:run"]
    },
    "prod": {
      "environment": {
        "SPRING_PROFILES_ACTIVE": "prod",
        "SERVER_ADDRESS": "0.0.0.0",
        "JAVA_OPTS": "-Xms512m -Xmx1024m"
      },
      "setup": ["./mvnw clean package -DskipTests"],
      "commands": ["java $JAVA_OPTS -jar target/*.jar"]
    },
    "port": 8080,
    "entrypoint": "#!/bin/bash\n\napp_env=${1:-development}\n\nif [ \"$app_env\" = \"production\" ] || [ \"$app_env\" = \"prod\" ] ; then\n    export SPRING_PROFILES_ACTIVE=prod\n    export SERVER_ADDRESS=0.0.0.0\n    export JAVA_OPTS=\"-Xms512m -Xmx1024m\"\n    ./mvnw clean package -DskipTests\n    java $JAVA_OPTS -jar target/*.jar\nelse\n    export SPRING_PROFILES_ACTIVE=dev\n    export SERVER_ADDRESS=0.0.0.0\n    ./mvnw spring-boot:run\nfi",
    "evidence": {
      "files": ["pom.xml", "src/main/java/", "mvnw"],
      "reason": "Java Spring Boot with Maven detected. Using mvnw wrapper."
    }
  }
]
```

**Key Points**:
- Maven detected from `pom.xml`
- Uses Maven wrapper (`mvnw`) for consistent builds
- Dev runs with Spring Boot Maven plugin for hot reload
- Production builds JAR and runs with optimized JVM settings
- Spring profiles control environment-specific config
- `SERVER_ADDRESS=0.0.0.0` ensures binding to all interfaces

---

## Example 2: Spring Boot with Gradle

**User request**: "Pack Gradle Spring Boot app"

**Detected files**:
- `build.gradle` or `build.gradle.kts`
- `gradlew` (Gradle wrapper)
- `src/main/java/`

**Output**:
```json
[
  {
    "language": "java",
    "version": "17",
    "apt": [],
    "dev": {
      "environment": {
        "SPRING_PROFILES_ACTIVE": "dev",
        "SERVER_ADDRESS": "0.0.0.0"
      },
      "commands": ["./gradlew bootRun"]
    },
    "prod": {
      "environment": {
        "SPRING_PROFILES_ACTIVE": "prod",
        "SERVER_ADDRESS": "0.0.0.0",
        "JAVA_OPTS": "-Xms512m -Xmx1024m"
      },
      "setup": ["./gradlew clean build -x test"],
      "commands": ["java $JAVA_OPTS -jar build/libs/*.jar"]
    },
    "port": 8080,
    "entrypoint": "#!/bin/bash\n\napp_env=${1:-development}\n\nif [ \"$app_env\" = \"production\" ] || [ \"$app_env\" = \"prod\" ] ; then\n    export SPRING_PROFILES_ACTIVE=prod\n    export SERVER_ADDRESS=0.0.0.0\n    export JAVA_OPTS=\"-Xms512m -Xmx1024m\"\n    ./gradlew clean build -x test\n    java $JAVA_OPTS -jar build/libs/*.jar\nelse\n    export SPRING_PROFILES_ACTIVE=dev\n    export SERVER_ADDRESS=0.0.0.0\n    ./gradlew bootRun\nfi",
    "evidence": {
      "files": ["build.gradle", "gradlew", "src/main/java/"],
      "reason": "Spring Boot with Gradle detected. Using Gradle wrapper."
    }
  }
]
```

**Key Points**:
- Gradle detected from `build.gradle` or `build.gradle.kts`
- Uses Gradle wrapper (`gradlew`)
- Dev: `bootRun` task with hot reload
- Prod: Build JAR to `build/libs/`, skip tests with `-x test`

---

## Example 3: Plain Java Application (Maven)

**User request**: "Analyze Maven project"

**Detected files**:
- `pom.xml` (no Spring Boot)
- `src/main/java/`
- Main class defined in pom.xml

**Output**:
```json
[
  {
    "language": "java",
    "version": "11",
    "apt": [],
    "dev": {
      "environment": {},
      "commands": ["./mvnw exec:java"]
    },
    "prod": {
      "environment": {},
      "setup": ["./mvnw clean package"],
      "commands": ["java -jar target/*.jar"]
    },
    "port": 8080,
    "entrypoint": "#!/bin/bash\n\napp_env=${1:-development}\n\nif [ \"$app_env\" = \"production\" ] || [ \"$app_env\" = \"prod\" ] ; then\n    ./mvnw clean package\n    java -jar target/*.jar\nelse\n    ./mvnw exec:java\nfi",
    "evidence": {
      "files": ["pom.xml", "src/main/java/"],
      "reason": "Plain Java application with Maven build system."
    }
  }
]
```

---

## Java Detection Patterns

### Required Files
- `pom.xml` (Maven) OR `build.gradle`/`build.gradle.kts` (Gradle)
- `src/main/java/` directory
- `*.java` files

### Build Tool Detection
| Build Tool | Signature Files | Priority |
|------------|----------------|----------|
| Maven | `pom.xml`, `mvnw`, `mvnw.cmd` | Higher |
| Gradle | `build.gradle`, `build.gradle.kts`, `gradlew` | Lower |

**If both exist**: Prefer Maven (more common in enterprise)

### Framework Detection
| Framework | Detection | Default Port |
|-----------|-----------|--------------|
| Spring Boot | `spring-boot-starter-*` in dependencies | 8080 |
| Micronaut | `micronaut-*` in dependencies | 8080 |
| Quarkus | `quarkus-*` in dependencies | 8080 |
| Vert.x | `vertx-*` in dependencies | 8080 |
| Plain Java | No framework dependencies | 8080 |

### Host Binding Patterns

**Spring Boot** (application.properties):
```properties
server.address=0.0.0.0
server.port=8080
```

**Via Environment Variable**:
```bash
export SERVER_ADDRESS=0.0.0.0
export SERVER_PORT=8080
```

**Programmatic** (in code):
```java
// Usually configured via properties
// Binding to 0.0.0.0 is default for most frameworks
```

### Version Detection
1. `pom.xml` → `<properties><java.version>` or `<maven.compiler.source>`
2. `build.gradle` → `sourceCompatibility` / `targetCompatibility`
3. `system.properties` file (Heroku-style)
4. Omit if uncertain

### Maven Commands

**Development**:
```bash
./mvnw spring-boot:run           # Spring Boot
./mvnw exec:java                 # Plain Java with exec plugin
./mvnw quarkus:dev               # Quarkus
```

**Production**:
```bash
./mvnw clean package             # Build JAR
./mvnw clean package -DskipTests # Skip tests
java -jar target/app.jar         # Run JAR
```

### Gradle Commands

**Development**:
```bash
./gradlew bootRun                # Spring Boot
./gradlew run                    # Plain Java
./gradlew quarkusDev             # Quarkus
```

**Production**:
```bash
./gradlew clean build            # Build JAR
./gradlew clean build -x test   # Skip tests
java -jar build/libs/app.jar    # Run JAR
```

### Common Environment Variables
- `SPRING_PROFILES_ACTIVE`: `dev` | `prod` | `test`
- `SERVER_ADDRESS`: `0.0.0.0` (binding address)
- `SERVER_PORT`: `8080` (port number)
- `JAVA_OPTS`: JVM options (e.g., `-Xms512m -Xmx1024m`)
- `JAVA_TOOL_OPTIONS`: Additional JVM options

### JVM Options for Production
```bash
# Memory settings
-Xms512m -Xmx1024m

# Garbage collection
-XX:+UseG1GC

# Performance monitoring
-XX:+PrintGCDetails -XX:+PrintGCDateStamps

# Heap dump on OOM
-XX:+HeapDumpOnOutOfMemoryError -XX:HeapDumpPath=/tmp
```

### APT Dependencies
**Default: `[]` (empty array)**

Reasons:
- JDK base images are complete and self-contained
- Java applications compile to bytecode (platform-independent)
- Dependencies managed by Maven/Gradle
- Native libraries bundled in JARs

**Only include APT packages if**:
- Project explicitly requires system libraries
- Native bindings needed (rare)
- Database client utilities: `postgresql-client`

### Output Locations
- **Maven**: JAR in `target/` directory
- **Gradle**: JAR in `build/libs/` directory
- **Fat JAR**: Usually named `*-SNAPSHOT.jar` or `*-all.jar`

### Spring Boot Specific

**Profiles**:
- `application-dev.properties` (dev)
- `application-prod.properties` (prod)
- Activated via `SPRING_PROFILES_ACTIVE`

**Configuration**:
```properties
# application.properties
server.address=0.0.0.0
server.port=8080
```

### Default Port
**8080** (standard for Java web applications)

### Production Best Practices
1. **JVM Tuning**: Set appropriate heap sizes with `JAVA_OPTS`
2. **Profiles**: Use Spring profiles for environment-specific config
3. **Fat JARs**: Use Spring Boot's executable JAR packaging
4. **Health Checks**: Enable actuator endpoints for monitoring
5. **Logging**: Configure proper logging levels for production
6. **Security**: Disable debug endpoints in production

### Wrapper Scripts
Always prefer wrappers over global Maven/Gradle:
- `./mvnw` instead of `mvn`
- `./gradlew` instead of `gradle`

Ensures consistent build tool versions across environments.
