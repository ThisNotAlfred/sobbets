project(
    'sobbet',
    'cpp',
    default_options: [
        'cpp_std=c++23',
        'warning_level=3',
    ],
)

sources = [
    'src/main.cpp',
]

includes = [
]

include_dirs = include_directories(
    [
        'include',
    ],
)

executable(
    'sobbet',
    [sources, includes],
    include_directories: [include_dirs],
)