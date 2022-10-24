#!/bin/bash

set -ex

dotfiles=/home/parigot/dotfiles

# Use this script after your host-side script has copied the appropriate files
# into this directory.

# docker cp leaves the user and group in a bogus state
sudo chown -R parigot ${dotfiles}
sudo chown -R parigot ${dotfiles}
sudo chmod -R 600 ${dotfiles}/ssh

if [ -d ${dotfiles}/ssh ]; then
  ln -s ${dotfiles}/ssh ${HOME}/.ssh
fi

if [ -d ${dotfiles}/oh-my-zsh ]; then
  ln -s ${dotfiles}/oh-my-zsh ${HOME}/.oh-my-zsh
fi

if [ -f ${dotfiles}/bashrc ]; then
  ln -s ${dotfiles}/bashrc ${HOME}/.bashrc
fi

if [ -f ${dotfiles}/bash_profile ]; then
  ln -s ${dotfiles}/bash_profile ${HOME}/.bash_profile
fi

if [ -f ${dotfiles}/zshrc ]; then
  ln -s ${dotfiles}/zshrc ${HOME}/.zshrc
fi

if [ -f ${dotfiles}/gitconfig ]; then
  ln -s ${dotfiles}/gitconfig ${HOME}/.gitconfig
fi
