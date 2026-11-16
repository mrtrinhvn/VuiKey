#----------------------------------------------------------------
# Generated CMake target import file for configuration "None".
#----------------------------------------------------------------

# Commands may need to know the format version.
set(CMAKE_IMPORT_FILE_VERSION 1)

# Import target "Fcitx5::Config" for configuration "None"
set_property(TARGET Fcitx5::Config APPEND PROPERTY IMPORTED_CONFIGURATIONS NONE)
set_target_properties(Fcitx5::Config PROPERTIES
  IMPORTED_LOCATION_NONE "${_IMPORT_PREFIX}/lib/x86_64-linux-gnu/libFcitx5Config.so.5.0.14"
  IMPORTED_SONAME_NONE "libFcitx5Config.so.6"
  )

list(APPEND _IMPORT_CHECK_TARGETS Fcitx5::Config )
list(APPEND _IMPORT_CHECK_FILES_FOR_Fcitx5::Config "${_IMPORT_PREFIX}/lib/x86_64-linux-gnu/libFcitx5Config.so.5.0.14" )

# Commands beyond this point should not need to know the version.
set(CMAKE_IMPORT_FILE_VERSION)
