## Will rebuild all ./templates in this subdirectory
## --supportBranch="blip/blipUtil" is needed because it is in this module instead of an included module
## projects using blip will not need to include this setting.
blip --supportBranch="blip/blipUtil" --rebuild
