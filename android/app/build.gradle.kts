plugins {
    id("com.android.application")
    id("org.jetbrains.kotlin.android")
}

android {
    namespace = "com.echipco.vuikey"
    compileSdk = 35

    defaultConfig {
        applicationId = "com.echipco.vuikey"
        minSdk = 24
        targetSdk = 36
        versionCode = 100
        versionName = "1.0.0"
        testInstrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"
    }

    signingConfigs {
        create("release") {
            val ksFile = file("keystore/vuikey-release.jks")
            if (ksFile.exists()) {
                storeFile = ksFile
                storePassword = System.getenv("VUIKEY_KS_PASS") ?: "vuikey2026"
                keyAlias = System.getenv("VUIKEY_KEY_ALIAS") ?: "vuikey"
                keyPassword = System.getenv("VUIKEY_KEY_PASS") ?: "vuikey2026"
            }
        }
    }

    buildTypes {
        release {
            isMinifyEnabled = true
            isShrinkResources = true
            val ksFile = file("keystore/vuikey-release.jks")
            signingConfig = if (ksFile.exists()) {
                signingConfigs.getByName("release")
            } else {
                signingConfigs.getByName("debug")
            }
            proguardFiles(
                getDefaultProguardFile("proguard-android-optimize.txt"),
                "proguard-rules.pro"
            )
        }
    }

    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_17
        targetCompatibility = JavaVersion.VERSION_17
    }

    kotlinOptions {
        jvmTarget = "17"
    }
}

dependencies {
    testImplementation("junit:junit:4.13.2")
}
