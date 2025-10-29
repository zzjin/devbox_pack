# Java Provider

- Detection: Uses confidence-based detection with weighted indicators:
  - Build files: `pom.xml`, `build.gradle`, `build.gradle.kts` (weight: 30)
  - Source files: `*.java`, `*.kt`, `*.scala` (weight: 25)
  - Gradle configuration: `gradle.properties` (weight: 15)
  - Gradle wrapper: `gradlew` (weight: 10)
  - Source directories: `src/main/java`, `src/main/kotlin`, `src/main/scala` (weight: 10)
  - Maven wrapper: `mvnw`, `mvnw.cmd` (weight: 5)
  - Build directories: `.mvn`, `.gradle` (weight: 5)
  - Detection threshold: 0.2 (20% confidence)

- Version Detection: Priority order for Java version resolution:
  1. `pom.xml` maven.compiler.source property
  2. `pom.xml` java.version property
  3. `pom.xml` source property
  4. `build.gradle` sourceCompatibility setting
  5. `build.gradle` targetCompatibility setting
  6. `build.gradle.kts` sourceCompatibility setting
  7. Default version: `17`

- Framework Detection: Automatically detects popular Java frameworks by analyzing build files:
  - **Spring**: Spring Boot, Spring MVC, Spring Framework
  - **Microservices**: Quarkus, Micronaut, Vert.x, Dropwizard
  - **Web**: Spark Java, Jersey, Struts, Apache Wicket, Vaadin
  - **Analysis**: Checks `pom.xml` and `build.gradle` for framework dependencies

- Package Manager Detection: Automatically detects build tool based on project files:
  - `pom.xml` → maven
  - `build.gradle` or `build.gradle.kts` → gradle
  - Default: maven

- Commands:
  - **Development**: 
    - **Maven**: `mvn clean compile`, then `mvn spring-boot:run`
    - **Gradle**: `./gradlew build`, then `./gradlew bootRun`
  - **Build**: 
    - **Maven**: `mvn clean package`
    - **Gradle**: `./gradlew build`
  - **Start**: 
    - **Maven**: `java -jar target/*.jar`
    - **Gradle**: `java -jar build/libs/*.jar`

- Native Compilation Detection: Java projects typically don't require native compilation
  - Returns `false` unless using GraalVM Native Image (not currently detected)

- Metadata: Provides comprehensive metadata including:
  - `hasPom`: Presence of `pom.xml`
  - `hasGradle`: Presence of `build.gradle` or `build.gradle.kts`
  - `hasGradleWrapper`: Presence of `gradlew`
  - `hasMavenWrapper`: Presence of `mvnw`
  - `packageManager`: Detected build tool (maven/gradle)
  - `framework`: Detected framework name