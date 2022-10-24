#!/bin/bash
# first and only parameters should be a container name (or id)

# you can use this to push things to a running dev container, the use the recv.sh
# script in ~/dotfiles to patch up permissions and such and symlink the files into
# the right places
set -ex

docker cp ~/.zshrc $1:/home/parigot/dotfiles/zshrc
docker cp ~/.gitconfig $1:/home/parigot/dotfiles/gitconfig
docker cp ~/.oh-my-zsh $1:/home/parigot/dotfiles/oh-my-zsh
docker cp ~/.ssh/config $1:/home/parigot/dotfiles/ssh/config
docker cp ~/.ssh/id_rsa $1:/home/parigot/dotfiles/ssh/id_rsa
docker cp ~/.ssh/id_rsa.pub $1:/home/parigot/dotfiles/ssh/id_rsa.pub

