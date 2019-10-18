# Breadcrumb trail towards APK build

## Install the Android SDK

`brew cask install android-sdk`

### Initialise config file

`touch ~/.android/repositories.cfg`

## Get a working version of Java

`brew cask install homebrew/cask-versions/adoptopenjdk8`

### Install Java from the package

`open /usr/local/Caskroom/adoptopenjdk8/8\,212\:b04/`

Double click the PKG, follow the prompts.

## Install the NDK bundle

`JAVA_HOME=/Library/Java/JavaVirtualMachines/adoptopenjdk-8.jdk/Contents/Home/ sdkmanager --install ndk-bundle`

## Build the APK

`ANDROID_NDK_HOME=/usr/local/Caskroom/android-sdk/4333796/ndk-bundle gomobile build -target=android golang.org/x/mobile/example/basic`

## References

- [Easily setup an Android development environment on a Mac](https://gist.github.com/patrickhammond/4ddbe49a67e5eb1b9c03)
- [Improve Android NDK path #58883](https://github.com/Homebrew/homebrew-cask/issues/58883)
- [Java8 not working anymore #7253](https://github.com/Homebrew/homebrew-cask-versions/issues/7253)
