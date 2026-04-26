# Project-level configuration.
cmake_minimum_required(VERSION 3.13)
project(runner LANGUAGES CXX)

# The name of the executable created for the application. Change this to change
# the on-disk name of your application.
set(BINARY_NAME "gen")
# The unique GTK application identifier for this application. See:
# https://wiki.gnome.org/HowDoI/ChooseApplicationID
set(APPLICATION_ID "ru.magomedcoder.gen")

# Explicitly opt in to modern CMake behaviors to avoid warnings with recent
# versions of CMake.
cmake_policy(SET CMP0063 NEW)

# Load bundled libraries from the lib/ directory relative to the binary.
set(CMAKE_INSTALL_RPATH "$ORIGIN/lib")

# Root filesystem for cross-building.
if(FLUTTER_TARGET_PLATFORM_SYSROOT)
  set(CMAKE_SYSROOT ${FLUTTER_TARGET_PLATFORM_SYSROOT})
  set(CMAKE_FIND_ROOT_PATH ${CMAKE_SYSROOT})
  set(CMAKE_FIND_ROOT_PATH_MODE_PROGRAM NEVER)
  set(CMAKE_FIND_ROOT_PATH_MODE_PACKAGE ONLY)
  set(CMAKE_FIND_ROOT_PATH_MODE_LIBRARY ONLY)
  set(CMAKE_FIND_ROOT_PATH_MODE_INCLUDE ONLY)
endif()

# Define build configuration options.
if(NOT CMAKE_BUILD_TYPE AND NOT CMAKE_CONFIGURATION_TYPES)
  set(CMAKE_BUILD_TYPE "Debug" CACHE
    STRING "Flutter build mode" FORCE)
  set_property(CACHE CMAKE_BUILD_TYPE PROPERTY STRINGS
    "Debug" "Profile" "Release")
endif()

# Compilation settings that should be applied to most targets.
#
# Be cautious about adding new options here, as plugins use this function by
# default. In most cases, you should add new options to specific targets instead
# of modifying this function.
function(APPLY_STANDARD_SETTINGS TARGET)
  target_compile_features(${TARGET} PUBLIC cxx_std_14)
  target_compile_options(${TARGET} PRIVATE -Wall -Werror)
  target_compile_options(${TARGET} PRIVATE "$<$<NOT:$<CONFIG:Debug>>:-O3>")
  target_compile_definitions(${TARGET} PRIVATE "$<$<NOT:$<CONFIG:Debug>>:NDEBUG>")
endfunction()

# Flutter library and tool build rules.
set(FLUTTER_MANAGED_DIR "${CMAKE_CURRENT_SOURCE_DIR}/flutter")
add_subdirectory(${FLUTTER_MANAGED_DIR})

# System-level dependencies.
find_package(PkgConfig REQUIRED)
pkg_check_modules(GTK REQUIRED IMPORTED_TARGET gtk+-3.0)

set(VOSK_BUNDLED_SO "${CMAKE_CURRENT_SOURCE_DIR}/native/vosk/libvosk.so")
set(VOSK_PREBUILT_VERSION "0.3.45" CACHE STRING "Версия ZIP vosk-linux-* на GitHub releases")

if(NOT DEFINED ENV{VOSK_LIB_FILE} AND NOT EXISTS "${VOSK_BUNDLED_SO}" AND NOT "$ENV{SKIP_VOSK_FETCH}" STREQUAL "1" AND NOT FLUTTER_TARGET_PLATFORM_SYSROOT)
  if(CMAKE_HOST_SYSTEM_PROCESSOR MATCHES "^(x86_64|AMD64)$")
    set(_VOSK_ARCH "x86_64")
  elseif(CMAKE_HOST_SYSTEM_PROCESSOR MATCHES "^(aarch64|ARM64)$")
    set(_VOSK_ARCH "aarch64")
  elseif(CMAKE_HOST_SYSTEM_PROCESSOR MATCHES "^(armv7l|armv6l)$")
    set(_VOSK_ARCH "armv7l")
  elseif(CMAKE_HOST_SYSTEM_PROCESSOR STREQUAL "riscv64")
    set(_VOSK_ARCH "riscv64")
  elseif(CMAKE_HOST_SYSTEM_PROCESSOR MATCHES "^(i.86|i686|x86)$")
    set(_VOSK_ARCH "x86")
  else()
    message(FATAL_ERROR
      "Vosk: неизвестная CMAKE_HOST_SYSTEM_PROCESSOR=${CMAKE_HOST_SYSTEM_PROCESSOR}. "
      "Положите libvosk.so в linux/native/vosk/ или задайте VOSK_LIB_FILE.")
  endif()

  set(_VOSK_NAME "vosk-linux-${_VOSK_ARCH}-${VOSK_PREBUILT_VERSION}")
  set(_VOSK_URL "https://github.com/alphacep/vosk-api/releases/download/v${VOSK_PREBUILT_VERSION}/${_VOSK_NAME}.zip")
  set(_VOSK_ZIP "${CMAKE_CURRENT_BINARY_DIR}/${_VOSK_NAME}.zip")
  set(_VOSK_EXTRACT "${CMAKE_CURRENT_BINARY_DIR}/vosk_prebuilt_extract")
  message(STATUS "Vosk: загрузка ${_VOSK_URL}")
  execute_process(
    COMMAND curl -fsSL -o "${_VOSK_ZIP}" "${_VOSK_URL}"
    RESULT_VARIABLE _VOSK_CURL
  )

  if(NOT _VOSK_CURL EQUAL 0)
    message(FATAL_ERROR "Vosk: curl не смог скачать архив")
  endif()

  file(REMOVE_RECURSE "${_VOSK_EXTRACT}")
  file(MAKE_DIRECTORY "${_VOSK_EXTRACT}")

  execute_process(
    COMMAND unzip -q -o "${_VOSK_ZIP}" -d "${_VOSK_EXTRACT}"
    RESULT_VARIABLE _VOSK_UNZIP
  )

  if(NOT _VOSK_UNZIP EQUAL 0)
    message(FATAL_ERROR "Vosk: unzip не смог распаковать архив.")
  endif()

  set(_VOSK_SO_SRC "${_VOSK_EXTRACT}/${_VOSK_NAME}/libvosk.so")

  if(NOT EXISTS "${_VOSK_SO_SRC}")
    message(FATAL_ERROR "Vosk: в архиве нет ${_VOSK_SO_SRC}")
  endif()

  file(MAKE_DIRECTORY "${CMAKE_CURRENT_SOURCE_DIR}/native/vosk")
  file(COPY "${_VOSK_SO_SRC}" DESTINATION "${CMAKE_CURRENT_SOURCE_DIR}/native/vosk")

  if(NOT EXISTS "${VOSK_BUNDLED_SO}")
    message(FATAL_ERROR "Vosk: не удалось установить ${VOSK_BUNDLED_SO}")
  endif()

  message(STATUS "Vosk: установлено ${VOSK_BUNDLED_SO}")
endif()

# Application build; see runner/CMakeLists.txt.
add_subdirectory("runner")

# Run the Flutter tool portions of the build. This must not be removed.
add_dependencies(${BINARY_NAME} flutter_assemble)

# Only the install-generated bundle's copy of the executable will launch
# correctly, since the resources must in the right relative locations. To avoid
# people trying to run the unbundled copy, put it in a subdirectory instead of
# the default top-level location.
set_target_properties(${BINARY_NAME}
  PROPERTIES
  RUNTIME_OUTPUT_DIRECTORY "${CMAKE_BINARY_DIR}/intermediates_do_not_run"
)


# Generated plugin build rules, which manage building the plugins and adding
# them to the application.
include(flutter/generated_plugins.cmake)


# === Installation ===
# By default, "installing" just makes a relocatable bundle in the build
# directory.
set(BUILD_BUNDLE_DIR "${PROJECT_BINARY_DIR}/bundle")
if(CMAKE_INSTALL_PREFIX_INITIALIZED_TO_DEFAULT)
  set(CMAKE_INSTALL_PREFIX "${BUILD_BUNDLE_DIR}" CACHE PATH "..." FORCE)
endif()

# Start with a clean build bundle directory every time.
install(CODE "
  file(REMOVE_RECURSE \"${BUILD_BUNDLE_DIR}/\")
  " COMPONENT Runtime)

set(INSTALL_BUNDLE_DATA_DIR "${CMAKE_INSTALL_PREFIX}/data")
set(INSTALL_BUNDLE_LIB_DIR "${CMAKE_INSTALL_PREFIX}/lib")

install(TARGETS ${BINARY_NAME} RUNTIME DESTINATION "${CMAKE_INSTALL_PREFIX}"
  COMPONENT Runtime)

install(FILES "${FLUTTER_ICU_DATA_FILE}" DESTINATION "${INSTALL_BUNDLE_DATA_DIR}"
  COMPONENT Runtime)

install(FILES "${FLUTTER_LIBRARY}" DESTINATION "${INSTALL_BUNDLE_LIB_DIR}"
  COMPONENT Runtime)

foreach(bundled_library ${PLUGIN_BUNDLED_LIBRARIES})
  install(FILES "${bundled_library}"
    DESTINATION "${INSTALL_BUNDLE_LIB_DIR}"
    COMPONENT Runtime)
endforeach(bundled_library)

# libvosk в bundle: VOSK_LIB_FILE или linux/native/vosk/libvosk.so (в т.ч. после авто-скачивания выше)
if(DEFINED ENV{VOSK_LIB_FILE})
  install(FILES "$ENV{VOSK_LIB_FILE}"
    DESTINATION "${INSTALL_BUNDLE_LIB_DIR}"
    COMPONENT Runtime)
elseif(EXISTS "${VOSK_BUNDLED_SO}")
  install(FILES "${VOSK_BUNDLED_SO}"
    DESTINATION "${INSTALL_BUNDLE_LIB_DIR}"
    COMPONENT Runtime)
endif()

# Copy the native assets provided by the build.dart from all packages.
set(NATIVE_ASSETS_DIR "${PROJECT_BUILD_DIR}native_assets/linux/")
install(DIRECTORY "${NATIVE_ASSETS_DIR}"
   DESTINATION "${INSTALL_BUNDLE_LIB_DIR}"
   COMPONENT Runtime)

# Fully re-copy the assets directory on each build to avoid having stale files
# from a previous install.
set(FLUTTER_ASSET_DIR_NAME "flutter_assets")
install(CODE "
  file(REMOVE_RECURSE \"${INSTALL_BUNDLE_DATA_DIR}/${FLUTTER_ASSET_DIR_NAME}\")
  " COMPONENT Runtime)
install(DIRECTORY "${PROJECT_BUILD_DIR}/${FLUTTER_ASSET_DIR_NAME}"
  DESTINATION "${INSTALL_BUNDLE_DATA_DIR}" COMPONENT Runtime)

# Install the AOT library on non-Debug builds only.
if(NOT CMAKE_BUILD_TYPE MATCHES "Debug")
  install(FILES "${AOT_LIBRARY}" DESTINATION "${INSTALL_BUNDLE_LIB_DIR}"
    COMPONENT Runtime)
endif()
