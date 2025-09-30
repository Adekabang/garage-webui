# Garage UI Style Guide

This document outlines the coding standards, naming conventions, and architectural patterns for the Garage UI project to ensure consistency across all development efforts.

## Table of Contents

- [Code Style](#code-style)
- [File and Folder Structure](#file-and-folder-structure)
- [Component Architecture](#component-architecture)
- [TypeScript Guidelines](#typescript-guidelines)
- [React Patterns](#react-patterns)
- [State Management](#state-management)
- [API and Data Fetching](#api-and-data-fetching)
- [Form Handling](#form-handling)
- [UI and Styling](#ui-and-styling)
- [Testing Guidelines](#testing-guidelines)
- [Git and Commit Guidelines](#git-and-commit-guidelines)

## Code Style

### General Principles

- **Consistency**: Follow established patterns in the codebase
- **Readability**: Write self-documenting code with clear naming
- **Simplicity**: Prefer simple, straightforward solutions
- **Performance**: Consider React performance best practices

### Formatting

- Use **2 spaces** for indentation
- Use **double quotes** for strings in JSX attributes
- Use **single quotes** for all other strings
- Use **semicolons** consistently
- Max line length: **100 characters**
- Use trailing commas in objects and arrays

```typescript
// ✅ Good
const config = {
  apiUrl: 'http://localhost:8080',
  timeout: 5000,
};

// ❌ Bad
const config = {
  apiUrl: "http://localhost:8080",
  timeout: 5000
}
```

### ESLint Configuration

Follow the existing ESLint configuration:
- TypeScript ESLint rules
- React Hooks rules
- React Refresh rules

## File and Folder Structure

### Naming Conventions

- **Files**: Use kebab-case for file names (`user-profile.tsx`, `auth-hooks.ts`)
- **Components**: Use PascalCase for component files (`UserProfile.tsx`, `NavigationBar.tsx`)
- **Folders**: Use kebab-case for folder names (`user-settings/`, `api-utils/`)
- **Assets**: Use kebab-case (`garage-logo.svg`, `user-avatar.png`)

### Folder Structure

```
src/
├── app/                    # App-level configuration
│   ├── app.tsx            # Main App component
│   ├── router.tsx         # Route definitions
│   ├── styles.css         # Global styles
│   └── themes.ts          # Theme configuration
├── assets/                # Static assets
├── components/            # Reusable components
│   ├── containers/        # Container components
│   ├── layouts/          # Layout components
│   └── ui/               # Basic UI components
├── context/               # React contexts
├── hooks/                 # Custom hooks
├── lib/                   # Utility libraries
├── pages/                 # Page components and related logic
│   └── [page-name]/
│       ├── index.tsx      # Main page component
│       ├── components/    # Page-specific components
│       ├── hooks.ts       # Page-specific hooks
│       ├── schema.ts      # Validation schemas
│       └── stores.ts      # Page-specific stores
├── stores/                # Global state stores
└── types/                 # TypeScript type definitions
```

### File Organization Rules

1. **Page Structure**: Each page should have its own folder with related components, hooks, schemas, and stores
2. **Component Isolation**: Page-specific components go in the page's `components/` folder
3. **Shared Components**: Reusable components go in `src/components/`
4. **Hooks**: Page-specific hooks in page folder, shared hooks in `src/hooks/`
5. **Types**: Domain-specific types in `src/types/`, component props types inline

## Component Architecture

### Component Types

1. **Page Components**: Top-level route components
2. **Layout Components**: Structural components (headers, sidebars, etc.)
3. **Container Components**: Components that manage state and logic
4. **UI Components**: Presentational components with minimal logic

### Component Structure

```typescript
// ✅ Good component structure
import { ComponentPropsWithoutRef, forwardRef } from 'react';
import { LucideIcon } from 'lucide-react';
import { Button as BaseButton } from 'react-daisyui';

// Types first
type ButtonProps = ComponentPropsWithoutRef<typeof BaseButton> & {
  icon?: LucideIcon;
  href?: string;
};

// Component with forwardRef for UI components
const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ icon: Icon, children, ...props }, ref) => {
    return (
      <BaseButton ref={ref} {...props}>
        {Icon && <Icon size={18} />}
        {children}
      </BaseButton>
    );
  }
);

Button.displayName = 'Button';

export default Button;
```

### Export Patterns

- **Default exports** for main components
- **Named exports** for utilities, hooks, and types
- **Barrel exports** for component directories (index.ts files)

```typescript
// utils.ts
export const formatDate = (date: Date) => { /* ... */ };
export const formatBytes = (bytes: number) => { /* ... */ };

// components/index.ts
export { default as Button } from './button';
export { default as Input } from './input';
```

## TypeScript Guidelines

### Type Definitions

- Use **interfaces** for object shapes that might be extended
- Use **types** for unions, primitives, and computed types
- Use **const assertions** for readonly arrays and objects

```typescript
// ✅ Good
interface User {
  id: string;
  name: string;
  email: string;
}

type Theme = 'light' | 'dark' | 'auto';

const themes = ['light', 'dark', 'auto'] as const;
type Theme = typeof themes[number];
```

### Generic Patterns

```typescript
// API response wrapper
type ApiResponse<T> = {
  data: T;
  success: boolean;
  message?: string;
};

// Component props with children
type ComponentProps<T = {}> = T & {
  children?: React.ReactNode;
  className?: string;
};
```

### Type Imports

Use type-only imports when importing only types:

```typescript
import type { User } from '@/types/user';
import type { ComponentProps } from 'react';
```

## React Patterns

### Hooks Usage

- **Custom hooks** for reusable logic
- **Built-in hooks** following React best practices
- **Hook naming**: Always start with `use`

```typescript
// ✅ Good custom hook
export const useAuth = () => {
  const { data, isLoading } = useQuery({
    queryKey: ['auth'],
    queryFn: () => api.get<AuthResponse>('/auth/status'),
    retry: false,
  });

  return {
    isLoading,
    isEnabled: data?.enabled,
    isAuthenticated: data?.authenticated,
  };
};
```

### Component Patterns

- Use **functional components** exclusively
- Use **forwardRef** for UI components that need ref access
- Destructure props in function parameters
- Use **early returns** for conditional rendering

```typescript
// ✅ Good component pattern
const UserCard = ({ user, onEdit, className }: UserCardProps) => {
  if (!user) {
    return <div className="text-gray-500">No user data</div>;
  }

  return (
    <div className={cn('card', className)}>
      <h3>{user.name}</h3>
      <p>{user.email}</p>
      <Button onClick={() => onEdit(user.id)}>Edit</Button>
    </div>
  );
};
```

## State Management

### Zustand Stores

- Use Zustand for **global state** management
- Keep stores **focused** and domain-specific
- Use **immer** for complex state updates

```typescript
// ✅ Good store pattern
import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';

interface AppState {
  theme: Theme;
  sidebarOpen: boolean;
  setTheme: (theme: Theme) => void;
  toggleSidebar: () => void;
}

export const useAppStore = create<AppState>()(
  immer((set) => ({
    theme: 'light',
    sidebarOpen: false,
    setTheme: (theme) => set((state) => {
      state.theme = theme;
    }),
    toggleSidebar: () => set((state) => {
      state.sidebarOpen = !state.sidebarOpen;
    }),
  }))
);
```

### Local State

- Use **useState** for simple local state
- Use **useReducer** for complex state logic
- Use **React Query** state for server state

## API and Data Fetching

### React Query Patterns

- Use **React Query** for all server state
- Follow consistent **query key** patterns
- Use **custom hooks** for API calls

```typescript
// ✅ Good API hook pattern
export const useUsers = (filters?: UserFilters) => {
  return useQuery({
    queryKey: ['users', filters],
    queryFn: () => api.get<User[]>('/users', { params: filters }),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

export const useCreateUser = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (userData: CreateUserInput) => 
      api.post<User>('/users', userData),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
    },
  });
};
```

### API Client

- Use consistent **error handling**
- Follow **RESTful** conventions
- Use **TypeScript** for request/response types

```typescript
// lib/api.ts
const api = {
  get: <T>(url: string, config?: AxiosRequestConfig) => 
    axios.get<T>(url, config).then(res => res.data),
  
  post: <T>(url: string, data?: any, config?: AxiosRequestConfig) => 
    axios.post<T>(url, data, config).then(res => res.data),
  
  // ... other methods
};
```

## Form Handling

### React Hook Form + Zod

- Use **React Hook Form** for all forms
- Use **Zod** for validation schemas
- Use **@hookform/resolvers** for integration

```typescript
// ✅ Good form pattern
// schema.ts
export const createUserSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  email: z.string().email('Invalid email'),
  role: z.enum(['admin', 'user']),
});

export type CreateUserInput = z.infer<typeof createUserSchema>;

// component.tsx
const CreateUserForm = ({ onSubmit }: CreateUserFormProps) => {
  const form = useForm<CreateUserInput>({
    resolver: zodResolver(createUserSchema),
    defaultValues: {
      name: '',
      email: '',
      role: 'user',
    },
  });

  return (
    <form onSubmit={form.handleSubmit(onSubmit)}>
      <Input
        {...form.register('name')}
        error={form.formState.errors.name?.message}
      />
      <Button type="submit" loading={form.formState.isSubmitting}>
        Create User
      </Button>
    </form>
  );
};
```

## UI and Styling

### TailwindCSS + DaisyUI

- Use **TailwindCSS** utility classes
- Use **DaisyUI** components as base
- Use **clsx** or **tailwind-merge** for conditional classes

```typescript
import { cn } from '@/lib/utils'; // tailwind-merge utility

const Button = ({ variant, size, className, ...props }: ButtonProps) => {
  return (
    <button
      className={cn(
        'btn', // DaisyUI base class
        {
          'btn-primary': variant === 'primary',
          'btn-secondary': variant === 'secondary',
          'btn-sm': size === 'small',
          'btn-lg': size === 'large',
        },
        className
      )}
      {...props}
    />
  );
};
```

### Theme System

- Use **DaisyUI themes** for consistent theming
- Define theme types for TypeScript support
- Use CSS custom properties for dynamic theming

```typescript
// themes.ts
export const themes = [
  'pastel',
  'dark',
  'dracula',
  // ... more themes
] as const;

export type Theme = typeof themes[number];
```

### Responsive Design

- **Mobile-first** approach
- Use Tailwind **responsive prefixes**
- Test on multiple screen sizes

```tsx
<div className="
  grid 
  grid-cols-1 
  md:grid-cols-2 
  lg:grid-cols-3 
  gap-4 
  p-4
">
  {/* Content */}
</div>
```

### Component Styling Guidelines

1. **Consistent spacing**: Use Tailwind spacing scale (4, 8, 12, 16, etc.)
2. **Color palette**: Use DaisyUI semantic colors (primary, secondary, accent, etc.)
3. **Typography**: Use Tailwind typography utilities
4. **Shadows and borders**: Use consistent elevation patterns

## Testing Guidelines

### Testing Strategy

1. **Unit tests** for utilities and pure functions
2. **Component tests** for UI components
3. **Integration tests** for complex workflows
4. **E2E tests** for critical user journeys

### Testing Patterns

```typescript
// ✅ Good test structure
describe('UserCard', () => {
  const mockUser = {
    id: '1',
    name: 'John Doe',
    email: 'john@example.com',
  };

  it('renders user information correctly', () => {
    render(<UserCard user={mockUser} />);
    
    expect(screen.getByText('John Doe')).toBeInTheDocument();
    expect(screen.getByText('john@example.com')).toBeInTheDocument();
  });

  it('calls onEdit when edit button is clicked', () => {
    const onEdit = jest.fn();
    render(<UserCard user={mockUser} onEdit={onEdit} />);
    
    fireEvent.click(screen.getByText('Edit'));
    expect(onEdit).toHaveBeenCalledWith('1');
  });
});
```

## Git and Commit Guidelines

### Branch Naming

- **feature/**: New features (`feature/user-authentication`)
- **fix/**: Bug fixes (`fix/login-validation`)
- **refactor/**: Code refactoring (`refactor/api-client`)
- **docs/**: Documentation (`docs/style-guide`)

### Commit Messages

Follow conventional commits:

```
feat: add user authentication system
fix: resolve login form validation issue
refactor: extract API client utilities
docs: update component documentation
style: fix linting issues
test: add user service tests
```

### PR Guidelines

1. **Clear description** of changes
2. **Link to issues** if applicable
3. **Screenshots** for UI changes
4. **Test coverage** for new features
5. **Breaking changes** clearly documented

---

## Quick Reference

### Import Order

1. React imports
2. Third-party libraries
3. Internal imports (components, hooks, utils)
4. Relative imports
5. Type-only imports (with `type` keyword)

### File Naming Quick Reference

| Type | Pattern | Example |
|------|---------|---------|
| Components | PascalCase.tsx | `UserProfile.tsx` |
| Hooks | camelCase.ts | `useAuth.ts` |
| Utils | kebab-case.ts | `api-client.ts` |
| Types | kebab-case.ts | `user-types.ts` |
| Pages | kebab-case.tsx | `user-settings.tsx` |

### Common Patterns

```typescript
// Component with props
type Props = {
  // prop definitions
};

const Component = ({ prop1, prop2 }: Props) => {
  // component logic
  return <div>{/* JSX */}</div>;
};

// Custom hook
export const useFeature = () => {
  // hook logic
  return { data, loading, error };
};

// API function
export const apiFunction = async (params: Params): Promise<Response> => {
  // API call
};
```

This style guide should be treated as a living document and updated as the project evolves and new patterns emerge.

# Go Backend Style Guide

This section outlines the coding standards and best practices for the Go backend of the Garage UI project.

## Table of Contents

- [Code Structure](#go-code-structure)
- [Naming Conventions](#go-naming-conventions)
- [Formatting](#go-formatting)
- [Error Handling](#go-error-handling)
- [Comments and Documentation](#go-comments-and-documentation)
- [API Design](#go-api-design)
- [Testing](#go-testing)
- [Performance Considerations](#go-performance-considerations)

## Go Code Structure

### Package Organization

- Use **meaningful package names** that reflect their purpose
- Follow Go's standard package layout:
  - `/cmd` - Main applications
  - `/pkg` - Library code that can be used by external applications
  - `/internal` - Private code not meant for external use

```
backend/
├── cmd/
│   └── server/         # Main application entrypoint
├── internal/
│   ├── api/            # API handlers
│   ├── middleware/     # HTTP middleware
│   ├── service/        # Business logic
│   ├── storage/        # Data storage/persistence
│   └── config/         # Configuration
├── pkg/
│   ├── models/         # Shared models/types
│   └── utils/          # Shared utilities
└── test/               # Test helpers
```

### File Organization

- Each file should focus on a single responsibility
- File names should clearly indicate their contents
- Group related files in appropriately named packages
- Limit file size (aim for under 500 lines where possible)

## Go Naming Conventions

### Package Names

- Use **short, lowercase, single-word** names
- Avoid underscores or mixedCaps
- Avoid plural forms (e.g., use `item` not `items`)

```go
// ✅ Good package names
package middleware
package router
package config
package auth

// ❌ Bad package names
package authHandler
package Middleware
package auth_handlers
```

### Variable Names

- Use **camelCase** for variable names
- Use **short but descriptive** names in local scopes
- Use **more descriptive names** for exported variables

```go
// ✅ Good
var userID string
var isAdmin bool
var httpClient *http.Client

// ❌ Bad
var UserID string       // Exported without need
var is_admin bool       // Using snake_case
var c *http.Client      // Too short for exported variable
```

### Function and Method Names

- Use **camelCase** for unexported functions
- Use **PascalCase** for exported functions
- Name methods with verbs or verb phrases

```go
// ✅ Good
func (s *Server) Start() error
func parseConfig(data []byte) (*Config, error)

// ❌ Bad
func (s *Server) server_start() error  // Using snake_case
func Config(data []byte) (*Config, error)  // Not a verb
```

### Interface Names

- Name interfaces based on the behavior they describe
- Often use the `-er` suffix for interfaces describing actions

```go
// ✅ Good
type Renderer interface {
    Render() []byte
}

type ConfigLoader interface {
    LoadConfig(path string) (*Config, error)
}
```

### Constant Names

- Use **PascalCase** for exported constants
- Use **camelCase** for package-level unexported constants
- Group related constants in const blocks

```go
// ✅ Good
const (
    // StatusActive represents an active user
    StatusActive = "active"
    
    // StatusInactive represents an inactive user
    StatusInactive = "inactive"
    
    maxRetries = 3  // Unexported constant
)
```

## Go Formatting

### Code Formatting

- Use `gofmt` or `goimports` to format code
- Use standard Go formatting:
  - Tabs for indentation
  - No trailing whitespace
  - No line length limit (but aim for ~100 characters for readability)

### Import Organization

- Group imports into blocks:
  1. Standard library
  2. Third-party packages
  3. Your project's packages

```go
import (
    "encoding/json"
    "fmt"
    "net/http"
    
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/joho/godotenv"
    
    "Adekabang/garage-webui/schema"
    "Adekabang/garage-webui/utils"
)
```

### Block Organization

- Order struct fields and methods logically, grouped by functionality
- Organize type definitions in a sensible order (important types first)
- Keep related code together

## Go Error Handling

### Error Creation

- Use `errors.New` for simple errors
- Use `fmt.Errorf` for formatted errors
- Use `%w` verb to wrap errors for additional context while preserving the original error

```go
// ✅ Good
if user == nil {
    return nil, errors.New("user not found")
}

if err := db.QueryRow(query, id).Scan(&user.ID, &user.Name); err != nil {
    return nil, fmt.Errorf("failed to query user %s: %w", id, err)
}
```

### Error Handling

- Check errors immediately after function calls that return them
- Don't use `_` to ignore errors without good reason
- Return errors to the caller rather than handling them at every level

```go
// ✅ Good
file, err := os.Open(filename)
if err != nil {
    return fmt.Errorf("opening config file: %w", err)
}
defer file.Close()
```

### Error Types

- Consider creating custom error types for specific error cases
- Use error wrapping to add context while preserving the original error
- Use `errors.Is` and `errors.As` for error checking with wrapped errors

```go
// ✅ Good
type NotFoundError struct {
    Resource string
    ID       string
}

func (e NotFoundError) Error() string {
    return fmt.Sprintf("%s with ID %s not found", e.Resource, e.ID)
}

// Usage
if errors.Is(err, NotFoundError{Resource: "User"}) {
    // Handle not found case
}
```

## Go Comments and Documentation

### Package Documentation

- Every package should have a package comment
- The package comment should be immediately before the `package` clause
- Describe the purpose of the package

```go
// Package middleware provides HTTP middleware functions for the Garage UI backend.
// It includes authentication, logging, and request validation middleware.
package middleware
```

### Function Documentation

- Document all exported functions, types, and variables
- Write comments in complete sentences
- Begin comments with the name of the thing being documented

```go
// GetUser retrieves a user by ID from the database.
// It returns a NotFoundError if the user does not exist.
func GetUser(id string) (*User, error) {
    // implementation
}
```

### Comment Style

- Use `//` for all comments, even multi-line comments
- Keep comments up to date with code changes
- Avoid obvious comments that don't add value

```go
// ✅ Good
// Calculate the average score from all submitted grades.
// Returns -1 if no grades have been submitted.
func calculateAverage(grades []int) float64 {
    // Implementation
}

// ❌ Bad
// This function calculates average
func calculateAverage(grades []int) float64 {
    // Implementation
}
```

## Go API Design

### Handler Structure

- Use the http.Handler interface
- Group related handlers in structs
- Use dependency injection for handler dependencies

```go
// ✅ Good pattern
type BucketHandler struct {
    garageClient GarageClient
    logger       Logger
}

func NewBucketHandler(client GarageClient, logger Logger) *BucketHandler {
    return &BucketHandler{
        garageClient: client,
        logger:       logger,
    }
}

func (h *BucketHandler) ListBuckets(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

### Middleware Pattern

- Use the standard middleware pattern
- Chain middleware in a clear order
- Keep middleware focused on a single responsibility

```go
// ✅ Good middleware pattern
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        duration := time.Since(start)
        log.Printf("%s %s %s", r.Method, r.URL.Path, duration)
    })
}
```

### Request/Response Handling

- Use structs for request and response data
- Validate request data early
- Return consistent error responses
- Use proper HTTP status codes

```go
// ✅ Good request handling
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        RespondWithError(w, http.StatusBadRequest, "Invalid request format")
        return
    }
    
    if err := validateUserRequest(req); err != nil {
        RespondWithError(w, http.StatusBadRequest, err.Error())
        return
    }
    
    user, err := h.service.CreateUser(req)
    if err != nil {
        RespondWithError(w, http.StatusInternalServerError, "Failed to create user")
        return
    }
    
    RespondWithJSON(w, http.StatusCreated, user)
}
```

## Go Testing

### Test Organization

- Place tests in the same package as the code being tested
- Use table-driven tests for multiple test cases
- Test both success and error cases

```go
// ✅ Good test structure
func TestUserService_GetUser(t *testing.T) {
    tests := []struct {
        name     string
        userID   string
        wantUser *User
        wantErr  bool
    }{
        {
            name:     "valid user",
            userID:   "123",
            wantUser: &User{ID: "123", Name: "Test User"},
            wantErr:  false,
        },
        {
            name:     "invalid user",
            userID:   "999",
            wantUser: nil,
            wantErr:  true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Mocking Dependencies

- Use interfaces to make testing easier
- Create mock implementations for testing
- Consider using a mocking library for complex interfaces

```go
// ✅ Good mocking pattern
type GarageClientMock struct {
    ListBucketsFn func() ([]Bucket, error)
}

func (m *GarageClientMock) ListBuckets() ([]Bucket, error) {
    return m.ListBucketsFn()
}

// In tests
func TestBucketHandler_ListBuckets(t *testing.T) {
    mockClient := &GarageClientMock{
        ListBucketsFn: func() ([]Bucket, error) {
            return []Bucket{{ID: "test"}}, nil
        },
    }
    
    handler := NewBucketHandler(mockClient, mockLogger)
    // Test implementation
}
```

### Test Helpers

- Create helper functions for common testing operations
- Use subtests for clearer test organization
- Make test failures descriptive

```go
// ✅ Good test helper pattern
func setupTest(t *testing.T) (*BucketHandler, *GarageClientMock, func()) {
    mockClient := &GarageClientMock{}
    mockLogger := &LoggerMock{}
    handler := NewBucketHandler(mockClient, mockLogger)
    
    return handler, mockClient, func() {
        // Cleanup code
    }
}
```

## Go Performance Considerations

### Resource Management

- Always close resources (files, database connections, etc.)
- Use `defer` for cleanup operations
- Consider resource pooling for expensive resources

```go
// ✅ Good resource management
func readConfig(path string) ([]byte, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    
    return io.ReadAll(file)
}
```

### Concurrency Patterns

- Use goroutines judiciously
- Use appropriate synchronization (mutex, channels, etc.)
- Avoid goroutine leaks

```go
// ✅ Good concurrency pattern
func processItems(items []Item) []Result {
    resultCh := make(chan Result, len(items))
    var wg sync.WaitGroup
    
    for _, item := range items {
        wg.Add(1)
        go func(i Item) {
            defer wg.Done()
            result := processItem(i)
            resultCh <- result
        }(item)
    }
    
    go func() {
        wg.Wait()
        close(resultCh)
    }()
    
    var results []Result
    for r := range resultCh {
        results = append(results, r)
    }
    
    return results
}
```

### Memory Management

- Preallocate slices and maps when size is known
- Be cautious with large allocations
- Watch for inadvertent memory retention

```go
// ✅ Good memory management
func processLargeFile(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()
    
    scanner := bufio.NewScanner(file)
    buffer := make([]byte, 64*1024)  // 64KB buffer
    scanner.Buffer(buffer, 1024*1024)  // Max token size 1MB
    
    for scanner.Scan() {
        line := scanner.Text()
        processLine(line)
    }
    
    return scanner.Err()
}
```

## Example Implementations

### API Handler

```go
// BucketHandler handles HTTP requests for bucket operations.
type BucketHandler struct {
    service  BucketService
    logger   Logger
}

// NewBucketHandler creates a new BucketHandler with the given dependencies.
func NewBucketHandler(service BucketService, logger Logger) *BucketHandler {
    return &BucketHandler{
        service: service,
        logger:  logger,
    }
}

// ListBuckets returns all buckets the user has access to.
func (h *BucketHandler) ListBuckets(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    buckets, err := h.service.ListBuckets(ctx)
    if err != nil {
        h.logger.Error("failed to list buckets", "error", err)
        RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve buckets")
        return
    }
    
    RespondWithJSON(w, http.StatusOK, buckets)
}
```

### Middleware

```go
// AuthMiddleware validates authentication for protected routes.
func AuthMiddleware(authService AuthService) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            token := r.Header.Get("Authorization")
            if token == "" {
                token = extractTokenFromCookie(r)
            }
            
            if token == "" {
                RespondWithError(w, http.StatusUnauthorized, "Authentication required")
                return
            }
            
            userID, err := authService.ValidateToken(token)
            if err != nil {
                RespondWithError(w, http.StatusUnauthorized, "Invalid authentication token")
                return
            }
            
            // Add user to context
            ctx := context.WithValue(r.Context(), UserIDKey, userID)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

### Service Layer

```go
// BucketService defines the interface for bucket operations.
type BucketService interface {
    ListBuckets(ctx context.Context) ([]Bucket, error)
    GetBucket(ctx context.Context, id string) (*Bucket, error)
    CreateBucket(ctx context.Context, bucket CreateBucketParams) (*Bucket, error)
    DeleteBucket(ctx context.Context, id string) error
}

// GarageBucketService implements BucketService for Garage.
type GarageBucketService struct {
    client GarageClient
}

// NewGarageBucketService creates a new GarageBucketService.
func NewGarageBucketService(client GarageClient) *GarageBucketService {
    return &GarageBucketService{
        client: client,
    }
}

// ListBuckets retrieves all buckets from Garage.
func (s *GarageBucketService) ListBuckets(ctx context.Context) ([]Bucket, error) {
    resp, err := s.client.ListBuckets(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to list buckets: %w", err)
    }
    
    var buckets []Bucket
    for _, b := range resp {
        buckets = append(buckets, mapToBucket(b))
    }
    
    return buckets, nil
}
```
