#!/bin/bash
#Build pam_tupass.so
USER=$(whoami)
gcc -fPIC -fno-stack-protector -c pam_tupass.c
ld -x --shared -o pam_tupass.so pam_tupass.o /usr/lib/libtupass.so

# Move neccessary files to their respective places for execution
if ! [ -x "$(command -v sudo)" ]; then
    if [ $USER = "root" ]; then
        [ -d /lib/security ] || mkdir /lib/security
        cp pam_tupass.so /lib/security/pam_tupass.so
        cp pam-config-tupass /usr/share/pam-configs/tupass
        echo -e "\e[31mRun 'sudo pam-auth-update' and enable TUPass PAM module.\e[0m"
    else
        echo "Error: sudo was not found and you are not root, please move files into appropriate locations"
        echo "pam_tupass.so -> /lib/security/pam_tupass.so and pam-config-tupass -> /usr/share/pam-configs/tupass"
        echo "Run 'pam-auth-update' to enable TUPass PAM module."
    fi
else
    [ -d /lib/security ] || sudo mkdir /lib/security
    sudo cp pam_tupass.so /lib/security/pam_tupass.so
    sudo cp pam-config-tupass /usr/share/pam-configs/tupass
    echo -e "\e[31mRun 'sudo pam-auth-update' and enable TUPass PAM module.\e[0m"
fi
