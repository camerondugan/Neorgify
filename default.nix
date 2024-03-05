{pkgs ? import <nixpkgs> 
  {
    config.allowUnfree = true; config.android_sdk.accept_license = true;
  }
}:

let
  buildToolsVersion = "34.0.5";
  cmakeVersion = "3.10.2";
  # Use buildToolsVersion when you define androidComposition
  androidComposition = pkgs.androidenv.composeAndroidPackages {
    # cmdLineToolsVersion = "8.0";
    # toolsVersion = "26.1.1";
    # platformToolsVersion = "34.0.5";
    # buildToolsVersions = [ "34.0.0" ];
    includeEmulator = true;
    # emulatorVersion = "34.1.9";
    platformVersions = [ "28" "29" "30" ];
    includeSources = false;
    includeSystemImages = false;
    systemImageTypes = [ "google_apis_playstore" ];
    # abiVersions = [ "armeabi-v7a" "arm64-v8a" ];
    cmakeVersions = [ cmakeVersion ];
    includeNDK = true;
    # ndkVersions = ["22.0.7026061"];
    # useGoogleAPIs = false;
    # useGoogleTVAddOns = false;
    # includeExtras = [
    #   "extras;google;gcm"
    # ];
  };
in
pkgs.mkShell rec {
  ANDROID_SDK_ROOT = "${androidComposition.androidsdk}/libexec/android-sdk";
  ANDROID_NDK_ROOT = "${ANDROID_SDK_ROOT}/ndk-bundle";

  # Use the same cmakeVersion here
  shellHook = ''
    export PATH="$(echo "$ANDROID_SDK_ROOT/cmake/${cmakeVersion}".*/bin):$PATH"
  '';

  # Use the same buildToolsVersion here
  # GRADLE_OPTS = "-Dorg.gradle.project.android.aapt2FromMavenOverride=${ANDROID_SDK_ROOT}/build-tools/${buildToolsVersion}/aapt2";

  buildInputs = with pkgs; [
    stdenv.cc.cc
    pkg-config
    libcxx
    libcxxabi
    libxml2
    xorg.libX11
    xorg.libXcursor
    xorg.libXfixes
    libxkbcommon
    vulkan-headers
    wayland
    glfw
  ];
}
