{ config, pkgs, ... }:

{
  imports = [ ./hardware-configuration.nix ];

  # Bootloader
  boot.loader.systemd-boot.enable = true;
  boot.loader.efi.canTouchEfiVariables = true;

  networking.hostName = "nixos-aaron";
  networking.networkmanager.enable = true;

  # Set your time zone and locale
  time.timeZone = "America/New_York"; # Change to your actual zone
  i18n.defaultLocale = "en_US.UTF-8";

  # --- GNOME & Windowing System ---
  services.xserver = {
    enable = true;
    displayManager.gdm.enable = true;
    desktopManager.gnome.enable = true;
    # Enable X11 support alongside Wayland
    layout = "us";
  };

  # Exclude unnecessary GNOME bloat
  environment.gnome.excludePackages = (with pkgs; [
    gnome-photos
    gnome-tour
    gedit
  ]) ++ (with pkgs.gnome; [
    cheese      # photo booth
    gnome-music
    gnome-terminal
    epiphany    # web browser
    geary       # email client
    evince      # document viewer
    gnome-characters
    totem       # video player
    tali        # poker game
    iagno       # go game
    hitori      # sudoku game
    atomix      # puzzle game
  ]);

  # --- Hardware & Audio ---
  services.xserver.libinput.enable = true; # Touchpad support
  
  security.rtkit.enable = true;
  services.pipewire = {
    enable = true;
    alsa.enable = true;
    alsa.support32Bit = true;
    pulse.enable = true;
    jack.enable = true;
  };

  # --- SSH & Networking ---
  services.openssh.enable = true;

  # --- User Account ---
  users.users.aaron = {
    isNormalUser = true;
    description = "Aaron";
    extraGroups = [ "networkmanager" "wheel" "video" "audio" ];
    # Password will be set during nixos-install or via 'passwd'
  };

  # --- System Packages ---
  environment.systemPackages = with pkgs; [
    # Essentials
    wget
    git
    micro
    pkg-config
    gnumake
    gcc
    binutils
    
    # Media & Apps
    code-cursor
    (mpv.override {
      scripts = [ mpvScripts.vapoursynth ];
    })
    vapoursynth
  ];

  # Allow unfree packages (needed for Cursor)
  nixpkgs.config.allowUnfree = true;

  system.stateVersion = "24.05"; # Check your current NixOS version
}