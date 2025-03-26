# How to change config relayer and resolver on enviroment

## Dependencies
- ansible-core

## Steps for change config
1. run decoding script
```
sh decoding_script.sh
```
2. Choose an option 'Decrypt configs'
3. Choose an environment
4. Choose a files (you can decrypt/encrypt all configs)
5. After decrypt you can edit configs in ./assets
6. Save changes 
7. Repeat steps from 2 to 4, but choose 'Encrypt configs'
8. Commit and push changes using git
