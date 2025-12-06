# TaskBoard Frontend

A modern, responsive task management application built with Next.js, TypeScript, and Tailwind CSS. Features a Kanban-style board with drag-and-drop functionality, JWT authentication, and real-time state management.

## Features

- 🔐 **JWT Authentication** - Secure signup/login with automatic token refresh
- 📋 **Kanban Board** - Visual task organization with To Do, In Progress, and Done columns
- 🎯 **Priority Indicators** - Color-coded borders (High=Red, Medium=Orange, Low=Blue)
- 🖱️ **Drag & Drop** - Intuitive task movement between columns using @dnd-kit
- ⚡ **Optimistic Updates** - Instant UI feedback with automatic rollback on errors
- 🌙 **Dark Mode Support** - Automatic theme detection
- 📱 **Responsive Design** - Works on desktop, tablet, and mobile

## Tech Stack

- **Framework**: Next.js 14 (App Router)
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **State Management**: Zustand
- **Drag & Drop**: @dnd-kit
- **Icons**: Lucide React
- **Date Handling**: date-fns

## Project Structure

```
src/
├── app/                    # Next.js App Router pages
│   ├── dashboard/          # Protected dashboard with task board
│   ├── login/              # Login page
│   ├── signup/             # Signup page
│   ├── globals.css         # Global styles
│   ├── layout.tsx          # Root layout
│   └── page.tsx            # Landing page (redirects)
├── components/
│   ├── auth/               # Authentication components
│   │   ├── LoginForm.tsx
│   │   └── SignupForm.tsx
│   ├── tasks/              # Task management components
│   │   ├── Header.tsx      # Dashboard header
│   │   ├── TaskBoard.tsx   # Main Kanban board
│   │   ├── TaskCard.tsx    # Individual task cards
│   │   ├── TaskColumn.tsx  # Status columns
│   │   └── TaskModal.tsx   # Task create/edit modal
│   └── ui/                 # Reusable UI components
│       ├── Button.tsx
│       ├── Input.tsx
│       ├── Modal.tsx
│       ├── Select.tsx
│       └── Textarea.tsx
├── hooks/
│   └── useAuth.ts          # Authentication hook
├── lib/
│   ├── api-client.ts       # API client with JWT handling
│   ├── config.ts           # API endpoints configuration
│   └── utils.ts            # Utility functions
├── services/
│   ├── auth-service.ts     # Auth API calls
│   └── task-service.ts     # Task API calls
├── stores/
│   ├── auth-store.ts       # Zustand auth state
│   └── task-store.ts       # Zustand task state
└── types/
    └── index.ts            # TypeScript type definitions
```

## Getting Started

### Prerequisites

- Node.js 18+ 
- npm or yarn

### Installation

1. Clone the repository:
   ```bash
   git clone <your-repo-url>
   cd task-management-frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Configure environment variables:
   ```bash
   cp .env.example .env.local
   ```
   
   Update `.env.local` with your Kong Gateway URL:
   ```
   NEXT_PUBLIC_API_BASE_URL=http://localhost:8000
   ```

4. Start the development server:
   ```bash
   npm run dev
   ```

5. Open [http://localhost:3000](http://localhost:3000) in your browser.

## API Configuration

Update `src/lib/config.ts` with your actual Kong Gateway endpoints if they differ:

```typescript
export const API_CONFIG = {
  BASE_URL: process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8000',
  
  AUTH: {
    SIGNUP: '/auth/signup',
    LOGIN: '/auth/login',
    LOGOUT: '/auth/logout',
    REFRESH: '/auth/refresh',
    VERIFY: '/auth/verify',
  },
  
  TASKS: {
    LIST: '/tasks',
    CREATE: '/tasks',
    GET: (id: string) => `/tasks/${id}`,
    UPDATE: (id: string) => `/tasks/${id}`,
    DELETE: (id: string) => `/tasks/${id}`,
    COMPLETE: (id: string) => `/tasks/${id}/complete`,
  },
};
```

## Backend API Requirements

### Auth Service Endpoints

| Method | Endpoint        | Description           |
|--------|----------------|-----------------------|
| POST   | /auth/signup   | Create new account    |
| POST   | /auth/login    | Login with credentials|
| POST   | /auth/logout   | Logout current user   |
| POST   | /auth/refresh  | Refresh access token  |
| GET    | /auth/verify   | Verify current token  |

### Task Service Endpoints

| Method | Endpoint             | Description          |
|--------|---------------------|----------------------|
| GET    | /tasks              | List all tasks       |
| POST   | /tasks              | Create new task      |
| GET    | /tasks/:id          | Get single task      |
| PUT    | /tasks/:id          | Update task          |
| DELETE | /tasks/:id          | Delete task          |
| PATCH  | /tasks/:id/complete | Mark task as done    |

## State Management

This project uses Zustand for state management, which provides:

- **Simple API** - Easy to understand and use
- **No boilerplate** - Minimal setup required  
- **TypeScript support** - Full type inference
- **DevTools support** - Easy debugging
- **Concurrent user support** - Each user session has isolated state

### Auth Store

Manages user authentication state including login, signup, logout, and token refresh.

### Task Store

Manages task CRUD operations with optimistic updates for better UX:
- Tasks are updated immediately in the UI
- If the API call fails, changes are automatically rolled back

## Customization

### Priority Colors

Update priority colors in `tailwind.config.ts`:

```typescript
colors: {
  priority: {
    high: '#ef4444',    // Red
    medium: '#f97316',  // Orange
    low: '#3b82f6',     // Blue
  },
}
```

### Status Colors

Update status column colors in `tailwind.config.ts`:

```typescript
colors: {
  status: {
    todo: '#6366f1',        // Indigo
    'in-progress': '#f59e0b', // Amber
    done: '#10b981',        // Emerald
  },
}
```

## Building for Production

```bash
npm run build
npm start
```

## Docker Support

You can add this Dockerfile to containerize the frontend:

```dockerfile
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM node:20-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static
COPY --from=builder /app/public ./public
EXPOSE 3000
CMD ["node", "server.js"]
```

## Future Enhancements

- [ ] User friends system
- [ ] Shared task boards
- [ ] Real-time collaboration (WebSockets)
- [ ] Task comments and attachments
- [ ] Task filtering and search
- [ ] Calendar view
- [ ] Task assignments

## License

MIT
