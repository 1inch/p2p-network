#!/bin/bash

# Function to ask for a password and store it in a variable
ask_password() {
    read -sp "Enter your Ansible-wault password: " password 
    echo
}

# Function to display the main menu
show_menu() {
    echo "Choose an option:"
    echo "1. Decrypt configs"
    echo "2. Encrypt configs"
    echo "3. Exit"
}

# Function to display the environment menu
show_menu_env() {
    echo "Choose an environment:"
    echo "1. Dev Environment"
}


# Function to display file selection from a directory
list_files() {
    local env_path=$1
    echo "Files in $env_path:"
    files=($(ls -A "$env_path"))  # Get list of files
    for i in "${!files[@]}"; do
        echo "$((i+1)). ${files[$i]}"
    done
    echo "$(( ${#files[@]} + 1 )). All files"
}

# Function to handle file selection and process it
handle_file_choice_decrypt() {
    local env_path=$1
    list_files "$env_path"
    read -p "Choose a file by number or select 'All files' option: " file_choice
    if (( file_choice > 0 && file_choice <= ${#files[@]} )); then
        selected_file="${files[$((file_choice-1))]}"
        echo "Processing file: $env_path/$selected_file"
        ask_password
        echo "$password" > vault_password
        ansible-vault decrypt "$env_path/$selected_file" --vault-password-file vault_password
        rm -f vault_password
    elif (( file_choice == ${#files[@]} + 1 )); then
        echo "Processing all files in $env_path..."
        ask_password
        echo "$password" > vault_password
        ansible-vault decrypt $env_path/*.yaml --vault-password-file vault_password
        rm -f vault_password    
    else
        echo "Invalid choice, please select a valid file number."
    fi
}

handle_file_choice_encrypt() {
    local env_path=$1
    list_files "$env_path"
    read -p "Choose a file by number or select 'All files' option: " file_choice

    if (( file_choice > 0 && file_choice <= ${#files[@]} )); then
        selected_file="${files[$((file_choice-1))]}"
        echo "Processing file: $env_path/$selected_file"
        ask_password
        echo "$password" > vault_password
        ansible-vault encrypt "$env_path/$selected_file" --vault-password-file vault_password
        rm -f vault_password
    elif (( file_choice == ${#files[@]} + 1 )); then
        echo "Processing all files in $env_path..."
        ask_password
        echo "$password" > vault_password
        ansible-vault encrypt $env_path/*.yaml --vault-password-file vault_password
        rm -f vault_password
    else
        echo "Invalid choice, please select a valid file number."
    fi
}


# Function to handle the choice of operation
handle_choice() {
    case $1 in
        1)
            show_menu_env
            read -p "Select environment (1, 2, 3): " env_choice
            case $env_choice in
                1)
                    echo "You selected Dev Environment."
                    handle_file_choice_decrypt "assets"
                    ;;
                *)
                    echo "Invalid environment selection."
                    exit 1
                    ;;
            esac
            ;;
        2)
            show_menu_env
            read -p "Select environment (1, 2, 3): " env_choice
            case $env_choice in
                1)
                    echo "You selected Dev Environment."
                    handle_file_choice_encrypt "assets"
                    ;;
                *)
                    echo "Invalid environment selection."
                    exit 1
                    ;;
            esac
            ;;
        3)
            exit 0
            ;;
        *)
            echo "Invalid choice, please select 1, 2, or 3."
            ;;
    esac
}

# Main program
show_menu
read -p "Enter your choice (1, 2, 3): " choice
handle_choice $choice