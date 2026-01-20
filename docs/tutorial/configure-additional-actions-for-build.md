# Image User Configuration Tutorial

This guide walks you through setting up login users for your target OS image using OS Image Composer.

## Prerequisites

- Linux environment
- OS Image Composer tool configured
- Basic understanding of YAML configuration

## Step 1: Understanding User Configuration

OS Image Composer supports two types of user password configuration:

1. **Plaintext passwords** (for development/testing only)
2. **Hashed passwords** (recommended for production)

## Step 2: Generate Password Hashes

For production environments, generate secure password hashes:

```bash
# Generate SHA-512 hash for a password
python3 -c "import crypt; print(crypt.crypt('your_password', crypt.mksalt(crypt.METHOD_SHA512)))"

# Alternative using openssl
openssl passwd -6 your_password

# Interactive password prompt (recommended)
python3 -c "import crypt, getpass; print(crypt.crypt(getpass.getpass(), crypt.mksalt(crypt.METHOD_SHA512)))"
```

**Security Note:** Never commit plaintext passwords to version control.

## Step 3: Configure Users in Your Template

Edit your OS Image Composer template YAML file to include user configurations:

```yaml
# Basic user configuration examples
systemConfig:
  ...
  users:
    # Development user with plaintext password (NOT for production)
    - name: user
      password: "user"  # Do not commit real plaintext passwords
      groups: ["sudo"]
      
    # Production user with hashed password
    - name: admin
      hash_algo: "sha512"
      password: "$6$qisZydr7DPWjCwDk$uiFDXvewTwAqs4H0gO7lRkmc5j2IUiuxSA8Yi.kjN9aLu4w3vysV80mD6C/0DvaBPLYCWU2fJwatYxVASJVL20"
      groups: ["sudo"]  
```

## Step 4: Common User Groups

### Common User Groups

When configuring users, assign only groups that exist in a minimal Linux OS installation. Common groups include:

- **`users`** – Standard user group (default for most user accounts)
- **`sudo`** – Sudo access group (for administrative privileges; may be called `wheel` on some distributions)
- **`adm`** – System monitoring and log access (present on many distributions)
- **`audio`** – Access to audio devices
- **`video`** – Access to video devices
- **`dialout`** – Access to serial ports

> **Note:** The availability and purpose of groups can vary by distribution. Avoid specifying groups like `docker`, `plugdev`, or `systemd-journal` unless you know they exist in your target OS.

## Step 5: Build Your OS Image

Run OS Image Composer to build your image with the configured users.

## Step 6: Test User Login

Test logging in with your configured users:

```bash
# Switch to a configured user
su - user

# Test sudo access
whoami

# Check user's groups
id
```

## Security Best Practices

1. **Never use plaintext passwords in production**
2. **Use strong, unique passwords for each user**
3. **Regularly rotate passwords**
4. **Assign minimal required group permissions**
5. **Remove or disable unused accounts**
6. **Consider using SSH keys instead of passwords**

## Troubleshooting

**Common Issues:**

1. **User cannot login:** Check password hash generation and syntax
2. **No sudo access:** Verify user is in `wheel` or `sudo` group
3. **Permission denied:** Check group assignments for required resources

**Debugging:**

```bash
# Check if user exists
id username

# Verify password hash
sudo cat /etc/shadow | grep username

# Check group membership
groups username
```

# Configuring Custom Commands During Image Build

## Overview

The OS Image Composer supports executing custom commands during the image build process through the `configurations` section in image template files. These commands are executed in a chroot environment after all packages have been installed, allowing you to customize the system configuration, create files, download resources, or perform any other setup tasks needed for your custom image.

## How It Works

Custom configurations are executed during the rootfs build phase, specifically after customer packages have been installed. The commands run in a chroot environment within the image being built, giving them full access to modify the target system.

## Configuration Structure

The `configurations` section must be placed under `systemConfig` in your image template YAML file, alongside other system configuration options:

```yaml
systemConfig:
  name: your-config-name
  description: "Your system configuration description"
  
  # Package installation happens first
  packages:
    - ubuntu-minimal
    - systemd
    - wget              # Required for download commands
    - curl              # Required for API calls
    # ... other packages
  
  # User configuration
  users:
    - name: user
      password: "user"
      groups: ["sudo"]
  
  # Custom configurations execute after packages are installed
  configurations:
    - cmd: "touch /etc/dummy01.txt"
    - cmd: "echo 'dlstreamer x86_64 ubuntu24 image' > /etc/dummy01.txt"
    - cmd: "wget --no-check-certificate -O /etc/validate.sh https://example.com/validate.sh"
```

## Complete Template Structure

Here's how the configurations section fits within a complete image template:

```yaml
image:
  name: your-image-name
  version: "1.0"

target:
  os: ubuntu
  dist: ubuntu24
  arch: x86_64
  imageType: raw

disk:
  name: your-disk-config
  # ... disk configuration

systemConfig:
  name: your-system-config
  description: "Custom system configuration"
  
  packages:
    - ubuntu-minimal
    - systemd-boot
    - openssh-server
    - wget                    # Essential for download operations
    - curl                    # For API interactions
    - systemd                 # For systemctl commands
    # ... other required packages
  
  kernel:
    version: "6.12"
    # ... kernel configuration
    
  users:
    - name: user
      password: "user"
      groups: ["sudo"]
  
  # Custom configurations - executed after all packages are installed
  configurations:
    - cmd: "mkdir -p /opt/myapp"
    - cmd: "touch /etc/dummy01.txt"
    - cmd: "echo 'Custom image marker' > /etc/dummy01.txt"
    - cmd: "wget --no-check-certificate -O /etc/validate.sh https://example.com/validate.sh"
    - cmd: "chmod +x /etc/validate.sh"
    - cmd: "systemctl enable ssh"
```

## Real-World Example

From the [`ubuntu24-mah.yml`](../image-templates/ubuntu24-mah.yml) template:

```yaml
systemConfig:
  
  packages:
    ...
    - wget                    # Required for wget commands below
      
  configurations:
    - cmd: "touch /etc/dummy01.txt"
    - cmd: "echo 'dlstreamer x86_64 ubuntu24 image' > /etc/dummy01.txt"
    - cmd: "wget --no-check-certificate -O /etc/validate.sh https://raw.githubusercontent.com/open-edge-platform/os-image-composer/main/validate.sh"
```

## Configuration Examples

### Basic File Operations

```yaml
systemConfig:
  ...    
  configurations:
    - cmd: "mkdir -p /opt/myapp/config"
    - cmd: "touch /etc/myconfig.conf"
    - cmd: "echo 'CustomApp=enabled' > /etc/myconfig.conf"
    - cmd: "chmod 644 /etc/myconfig.conf"
```

### Download and Install Resources

```yaml
systemConfig:
  packages:
    ....
    - wget                 # Required for wget commands
    - curl                 # Required for curl commands
    
  configurations:
    - cmd: "wget --no-check-certificate -O /tmp/setup.sh https://example.com/setup.sh"
    - cmd: "chmod +x /tmp/setup.sh"
    - cmd: "bash /tmp/setup.sh"
    - cmd: "curl -o /etc/app.json https://api.example.com/config"
```

## Execution Environment

### Chroot Context

All commands are executed using `chroot` to the image root filesystem. This means:

- Commands run as if the target system is the root filesystem
- All paths are relative to the image being built
- Standard system tools and shell are available
- Network access is available (if needed for downloads)

### Execution Order

The build process follows this sequence:

1. **Package Installation**: All packages listed in `packages` are installed
2. **Custom Configurations**: Commands in `configurations` are executed sequentially


## Best Practices

### 1. Include Required Tools in Packages

Always include the tools your commands need in the `packages` section:

```yaml
systemConfig:
  packages:
    - wget              # For wget commands
    - curl              # For curl commands  
    - systemd           # For systemctl commands
    - coreutils         # For basic file operations
    - util-linux        # For system utilities
    
  configurations:
    - cmd: "wget -O /opt/file.txt https://example.com/file.txt"
    - cmd: "systemctl enable myservice"
```

### 2. Use Absolute Paths

Always use absolute paths since commands execute in the chroot environment:

```yaml
configurations:
  # Good - absolute paths
  - cmd: "echo 'config' > /etc/myapp.conf"
  - cmd: "mkdir -p /opt/myapp"
  
  # Avoid - relative paths may not work as expected
  - cmd: "echo 'config' > myapp.conf"
```

### 3. Make Commands Robust

Add error checking and make commands idempotent:

```yaml
configurations:
  - cmd: "mkdir -p /opt/myapp || true"
  - cmd: "test -f /etc/config || echo 'default' > /etc/config"
  - cmd: "systemctl is-enabled ssh || systemctl enable ssh"
```

### 4. Secure Downloads

Be cautious with downloaded content and verify integrity when possible:

```yaml
systemConfig:
  packages:
    - wget
    - coreutils         # For sha256sum
    
  configurations:
    - cmd: "wget --no-check-certificate -O /tmp/script.sh https://example.com/script.sh"
    - cmd: "echo 'expected-sha256 /tmp/script.sh' | sha256sum -c"
    - cmd: "chmod +x /tmp/script.sh && mv /tmp/script.sh /opt/script.sh"
```

## Error Handling

- If any configuration command fails, the entire build process stops
- Error messages are logged with details about which command failed
- Commands should be designed to be idempotent when possible

## Troubleshooting

### Common Issues

1. **Command Not Found**: Ensure required packages are installed in the `packages` section
2. **Permission Denied**: Some operations may require specific user contexts
3. **Network Failures**: Download commands may fail due to network issues
4. **Path Issues**: Use absolute paths for all file operations

## Security Considerations

- Downloaded scripts and files should be from trusted sources
- Avoid hardcoding sensitive information in commands
- Consider using checksums to verify downloaded content
- Be cautious with commands that modify system security settings

## Related Documentation

- [Image Template Format](../architecture/image-template-format.md)
- [Understanding the OS Image Build Process](../architecture/os-image-composer-build-process.md)
- [Command-Line Reference](../architecture/os-image-composer-cli-specification.md)