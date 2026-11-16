#----------------------------------------------------------------
# Generated CMake target import file for configuration "None".
#----------------------------------------------------------------

# Commands may need to know the format version.
set(CMAKE_IMPORT_FILE_VERSION 1)

# Import target "Fcitx5::Utils" for configuration "None"
set_property(TARGET Fcitx5::Utils APPEND PROPERTY IMPORTED_CONFIGURATIONS NONE)
set_target_properties(Fcitx5::Utils PROPERTIES
  IMPORTED_LOCATION_NONE "${_IMPORT_PREFIX}/lib/x86_64-linux-gnu/libFcitx5Utils.so.5.0.14"
  IMPORTED_SONAME_NONE "libFcitx5Utils.so.2"
  )

list(APPEND _IMPORT_CHECK_TARGETS Fcitx5::Utils )
list(APPEND _IMPORT_CHECK_FILES_FOR_Fcitx5::Utils "${_IMPORT_PREFIX}/lib/x86_64-linux-gnu/libFcitx5Utils.so.5.0.14" )

# Commands beyond this point should not need to know the version.
set(CMAKE_IMPORT_FILE_VERSION)
