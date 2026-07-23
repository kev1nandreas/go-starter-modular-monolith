#!/bin/bash

ENTITY_TYPE="$1"

show_usage() {
    echo "Usage: $0 [OPTION] [ENTITY_NAME]"
    echo ""
    echo "Options:"
    echo "  all         Process all items"
    echo "  <entity>    Process specific entity"
    echo "  help        Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 books"
    echo "  $0 users"
}

remove_lines_from_file() {
    local file="$1"
    local pattern="$2"
    
    if [ -f "$file" ]; then
        echo "  Removing lines matching '$pattern' from $file"
        # Create backup
        cp "$file" "${file}.bak"
        # Remove lines containing the pattern
        grep -v "$pattern" "${file}.bak" > "$file"
        rm "${file}.bak"
    else
        echo "  Warning: File $file not found"
    fi
}

remove_block_from_file() {
    local file="$1"
    local start_pattern="$2"
    local end_pattern="$3"
    
    if [ -f "$file" ]; then
        echo "  Removing block from $file"
        # Create backup
        cp "$file" "${file}.bak"
        # Remove lines between patterns (inclusive)
        sed "/$start_pattern/,/$end_pattern/d" "${file}.bak" > "$file"
        rm "${file}.bak"
    else
        echo "  Warning: File $file not found"
    fi
}

remove_route_block() {
    local file="pkg/api/router.go"
    local entity="$1"
    
    if [ -f "$file" ]; then
        echo "  Removing $entity routes from $file"
        cp "$file" "${file}.bak"
        
        # Use awk to remove the route block
        awk -v entity="$entity" '
        BEGIN { in_block = 0; skip = 0 }
        {
            # Check if line contains entity route (GET, POST, PUT, DELETE)
            if ($0 ~ "v1\\.(GET|POST|PUT|DELETE)\\(\"/?" entity) {
                skip = 1
                next
            }
            # If not skipping, print the line
            if (!skip) {
                print $0
            } else {
                skip = 0
            }
        }
        ' "${file}.bak" > "$file"
        
        rm "${file}.bak"
    else
        echo "  Warning: File $file not found"
    fi
}

process_entity() {
    local entity="$1"
    local entity_lower=$(echo "$entity" | tr '[:upper:]' '[:lower:]')
    local entity_title=$(echo "$entity" | awk '{print toupper(substr($0,1,1)) tolower(substr($0,2))}')
    
    echo "Processing $entity entity cleanup..."
    echo ""
    
    # Remove API files
    echo "1. Removing API files..."
    rm -f "pkg/api/${entity_lower}_mock.go"
    rm -f "pkg/api/${entity_lower}_test.go"
    rm -f "pkg/api/${entity_lower}.go"
    echo "  Removed pkg/api/${entity_lower}*.go files"
    echo ""
    
    # Remove repository initialization from router.go
    echo "2. Cleaning router.go..."
    remove_lines_from_file "pkg/api/router.go" "${entity_lower}Repository"
    remove_route_block "$entity_lower"
    echo ""
    
    # Remove model from migration
    echo "3. Cleaning pkg/database/migration/migrate.go..."
    remove_lines_from_file "pkg/database/migration/migrate.go" "\&models\.${entity_title}{"
    echo ""
    
    # Remove seeder calls
    echo "4. Cleaning pkg/database/seeders/seeder.go..."
    remove_lines_from_file "pkg/database/seeders/seeder.go" "Seed${entity_title}"
    remove_lines_from_file "pkg/database/seeders/seeder.go" "Clear${entity_title}"
    echo ""
    
    # Remove seeder file
    echo "5. Removing seeder file..."
    rm -f "pkg/database/seeders/${entity_lower}_seeder.go"
    echo "  Removed pkg/database/seeders/${entity_lower}_seeder.go"
    echo ""
    
    # Remove model file
    echo "6. Removing model file..."
    rm -f "pkg/models/${entity_lower}.go"
    echo "  Removed pkg/models/${entity_lower}.go"
}

process_all() {
    echo "Processing all items cleanup..."
    echo "This will remove all generated entity files."
    echo ""
    read -p "Are you sure? (y/n): " confirm
    if [ "$confirm" != "y" ]; then
        echo "Cancelled."
        exit 0
    fi
    
    # You can add multiple entities here
    process_entity "Book"
    process_entity "book"
    process_entity "books"
    process_entity "book"
}

if [ $# -eq 0 ]; then
    echo "Error: No arguments provided"
    show_usage
    exit 1
fi

# Parse command line arguments
case "$1" in
    all)
        process_all
        ;;
    help)
        show_usage
        exit 0
        ;;
    *)
        # Treat as entity name
        process_entity "$1"
        ;;
esac

echo ""
echo "Done!"