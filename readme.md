It's simple application for compare excel files places on different directories.

### How usage

1. Place default excel files to directory one.
2. Place new excel files to directory two.
3. Open terminal on your computer (on windows button "win + R", enter "cmd" on input string and press Enter button)4. Enter full path to application *.exe file (Default file name "locomparer.exe")

### Existing application flags:

    config - [not required] Full path to *.json configuration file (Default path to configuration file it directory
    when place application)
    dDir - [required] Full path to directory with default excel files
    nDir - [required] Full path to directory with new excel files
    oDir - Full path to directory for save compare results  (Default path it directory when place application)

Run application example:

    locomparer -dDir="C:\apps\def_test" -nDir="C:\apps\new_test" -oFile="C:\apps\compare.xlsx"
        -oDir="C:\apps\compare_xlsx"