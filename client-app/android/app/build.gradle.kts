import java.net.URI
import java.util.zip.ZipFile

plugins {
    id("com.android.application")
    id("kotlin-android")
    // The Flutter Gradle Plugin must be applied after the Android and Kotlin Gradle plugins.
    id("dev.flutter.flutter-gradle-plugin")
}

android {
    namespace = "ru.magomedcoder.gen"
    compileSdk = flutter.compileSdkVersion
    ndkVersion = flutter.ndkVersion

    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_17
        targetCompatibility = JavaVersion.VERSION_17
    }

    kotlinOptions {
        jvmTarget = JavaVersion.VERSION_17.toString()
    }

    defaultConfig {
        // TODO: Specify your own unique Application ID (https://developer.android.com/studio/build/application-id.html).
        applicationId = "ru.magomedcoder.gen"
        // You can update the following values to match your application needs.
        // For more information, see: https://flutter.dev/to/review-gradle-config.
        minSdk = flutter.minSdkVersion
        targetSdk = flutter.targetSdkVersion
        versionCode = flutter.versionCode
        versionName = flutter.versionName
    }

    buildTypes {
        release {
            // TODO: Add your own signing config for the release build.
            // Signing with the debug keys for now, so `flutter run --release` works.
            signingConfig = signingConfigs.getByName("debug")
        }
    }
}

flutter {
    source = "../.."
}

val voskPrebuiltVersion = providers.gradleProperty("vosk.prebuilt.version").getOrElse("0.3.45")

tasks.register("fetchVoskJniLibs") {
    group = "vosk"
    description = "Скачивает vosk-android-*.zip и раскладывает libvosk.so по jniLibs/*"
    val jniBase = file("src/main/jniLibs")
    val marker = jniBase.resolve("arm64-v8a/libvosk.so")
    outputs.file(marker)
    onlyIf {
        System.getenv("SKIP_VOSK_FETCH") != "1" && !marker.exists()
    }
    doLast {
        val ver = voskPrebuiltVersion
        val name = "vosk-android-$ver.zip"
        val url = URI("https://github.com/alphacep/vosk-api/releases/download/v$ver/$name").toURL()
        val zipFile = layout.buildDirectory.file(name).get().asFile
        zipFile.parentFile?.mkdirs()
        logger.lifecycle("Vosk Android: загрузка $url")
        url.openStream().use { input ->
            zipFile.outputStream().use { out -> input.copyTo(out) }
        }
        ZipFile(zipFile).use { zip ->
            zip.entries().asSequence().forEach { e ->
                if (e.isDirectory) {
                    return@forEach
                }
                val entryName = e.name.replace('\\', '/')
                if (!entryName.endsWith("libvosk.so")) {
                    return@forEach
                }
                val parent = entryName.substringBeforeLast('/', "")
                if (parent.isEmpty() || parent.contains("..")) {
                    return@forEach
                }
                val dest = jniBase.resolve("$parent/libvosk.so")
                dest.parentFile?.mkdirs()
                zip.getInputStream(e).use { input ->
                    dest.outputStream().use { out -> input.copyTo(out) }
                }
            }
        }
        logger.lifecycle("Vosk Android: jniLibs обновлены ($ver)")
    }
}

tasks.named("preBuild").configure { dependsOn("fetchVoskJniLibs") }
