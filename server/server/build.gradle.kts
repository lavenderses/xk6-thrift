plugins {
    alias(libs.plugins.kotlin.jvm)
    application
}

repositories {
    mavenCentral()
}

dependencies {
    implementation(project(":thrift"))

    implementation(platform(libs.armeria.bom))
    implementation(libs.armeria)
    implementation(libs.armeria.logback)
    implementation(libs.armeria.thrift)
    implementation(libs.logback.classic)

    implementation(libs.coroutine.core)

    testImplementation(libs.junit.jupiter)
    testRuntimeOnly("org.junit.platform:junit-platform-launcher")
}

java {
    toolchain {
        languageVersion = JavaLanguageVersion.of(21)
    }
}

application {
    mainClass = "server.AppKt"
}

tasks.named<Test>("test") {
    useJUnitPlatform()
}
