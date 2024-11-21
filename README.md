# Duet

Duet is a powerful infrastructure and configuration management tool that orchestrates both infrastructure provisioning and system configuration in harmony. Using Lua as its configuration language, Duet provides a unified approach to managing your entire infrastructure lifecycle.

## Features

- Infrastructure as Code (similar to Terraform/Pulumi)
- Configuration Management (similar to Ansible)
- Lua-based configuration language
- State management using SQLite
- Idempotent operations
- AWS provider support (more coming soon)

## Quick Start

```bash
# Install Duet
go install github.com/rebelopsio/duet/cmd/duet@latest

# Create a configuration file
cat > infra.lua << EOF
local config = {
    infrastructure = {
        aws = {
            region = "us-west-2",
            ec2 = {
                instance_type = "t2.micro",
                ami = "ami-0c55b159cbfafe1f0",
                subnet_id = "subnet-xxxxxxxx"
            }
        }
    }
}

function deploy_infrastructure()
    return config.infrastructure
end

function configure_instance(host)
    local success, err = install_package("cowsay", host)
    if not success then
        error("Failed to install cowsay: " .. err)
    end
    return true
end
EOF

# Plan your changes
duet plan infra.lua

# Apply your changes
duet apply infra.lua
```

## Architecture

Duet is built with a clear separation of concerns:

- Infrastructure as Code (IaC) Engine: Manages infrastructure provisioning
- Configuration Management Engine: Handles system configuration
- Lua Engine: Processes configuration files
- State Store: Maintains system state using SQLite

## Development

```bash
# Clone the repository
git clone https://github.com/rebelopsio/duet.git

# Install dependencies
go mod download

# Build
go build -o duet cmd/duet/main.go

# Run tests
go test ./...
```

## Contributing

Contributions are welcome! Please read our contributing guidelines before submitting pull requests.

## License

MIT License - see LICENSE file for details
