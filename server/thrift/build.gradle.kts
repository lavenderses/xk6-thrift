plugins {
    alias(libs.plugins.thrift.gradle.plugin)
    alias(libs.plugins.kotlin.jvm)
    idea
}

repositories {
    mavenCentral()
}

dependencies {
    implementation(libs.coroutine.core)
    implementation(libs.thrift)
    implementation(libs.slf4j.api)
    implementation(libs.jakarta.annotation.api)
}

java {
    toolchain {
        languageVersion = JavaLanguageVersion.of(21)
    }
}

tasks.named<com.linecorp.thrift.plugin.CompileThrift>("compileThrift") {
    recurse = true
    thriftExecutable = "${rootProject.file("bin/thrift")}"

    generators = mapOf(
        "java" to "",
        "json" to "",
        "html" to "",
        "go" to "",
    )
}

val thriftGeneratedSource = file(layout.buildDirectory.dir("generated-sources/gen-java"))

sourceSets.main {
    kotlin {
        srcDir(thriftGeneratedSource)
    }
}

idea {
    module {
        generatedSourceDirs.add(thriftGeneratedSource)
    }
}
