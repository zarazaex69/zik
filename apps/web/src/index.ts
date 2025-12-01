import { serve } from "bun";
import { readFileSync } from "fs";
import { join } from "path";
import index from "./index.html";

const server = serve({
  port: 8805,
  routes: {
    // Serve install script
    "/install": {
      async GET() {
        const scriptPath = join(import.meta.dir, "install.sh");
        const script = readFileSync(scriptPath, "utf-8");
        return new Response(script, {
          headers: {
            "Content-Type": "text/plain",
          },
        });
      },
    },

    // Serve favicon
    "/assets/favicon.ico": {
      async GET() {
        const faviconPath = join(import.meta.dir, "../assets/favicon.ico");
        const favicon = readFileSync(faviconPath);
        return new Response(favicon, {
          headers: {
            "Content-Type": "image/x-icon",
          },
        });
      },
    },

    // Serve index.html for all unmatched routes.
    "/*": index,

    "/api/hello": {
      async GET(req) {
        return Response.json({
          message: "Hello, world!",
          method: "GET",
        });
      },
      async PUT(req) {
        return Response.json({
          message: "Hello, world!",
          method: "PUT",
        });
      },
    },

    "/api/hello/:name": async req => {
      const name = req.params.name;
      return Response.json({
        message: `Hello, ${name}!`,
      });
    },
  },

  development: process.env.NODE_ENV !== "production" && {
    // Enable browser hot reloading in development
    hmr: true,

    // Echo console logs from the browser to the server
    console: true,
  },
});

console.log(`Server running at ${server.url}`);

// Graceful shutdown handler
const shutdown = async (signal: string) => {
  console.log(`\n${signal} received, shutting down gracefully...`);
  try {
    server.stop();
    console.log("Server stopped successfully");
    process.exit(0);
  } catch (error) {
    console.error("Error during shutdown:", error);
    process.exit(1);
  }
};

// Handle shutdown signals
process.on("SIGTERM", () => shutdown("SIGTERM"));
process.on("SIGINT", () => shutdown("SIGINT"));
