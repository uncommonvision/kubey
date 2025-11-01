# Agent Instructions for Kubey Web

## Commands
- **Build**: `npm run build` (runs TypeScript check + Vite build)
- **Dev server**: `npm run dev` (Vite dev server)
- **Lint**: `npm run lint` (ESLint on all files)
- **Preview**: `npm run preview` (Vite preview)

## Code Style
- **Framework**: React 19 + TypeScript + Vite
- **UI**: shadcn/ui with "new-york" style, Tailwind CSS, CSS variables
- **Imports**: Use `@/*` alias for src/, React imports first
- **Components**: forwardRef with displayName, cn() for class merging
- **Types**: Strict TypeScript, no unused locals/parameters
- **Linting**: ESLint recommended + React hooks + React refresh
- **Naming**: PascalCase for components, camelCase for functions/hooks
- **Error handling**: Use try/catch for async operations
- **Styling**: Tailwind classes, responsive design, semantic HTML