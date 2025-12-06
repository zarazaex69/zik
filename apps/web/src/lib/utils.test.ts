import { describe, test, expect } from "bun:test";
import { cn } from "./utils";

describe("cn utility function", () => {
  test("should merge single class name", () => {
    const result = cn("text-red-500");
    expect(result).toBe("text-red-500");
  });

  test("should merge multiple class names", () => {
    const result = cn("text-red-500", "bg-blue-500");
    expect(result).toContain("text-red-500");
    expect(result).toContain("bg-blue-500");
  });

  test("should handle conditional classes", () => {
    const isActive = true;
    const result = cn("base-class", isActive && "active-class");
    expect(result).toContain("base-class");
    expect(result).toContain("active-class");
  });

  test("should filter out false values", () => {
    const result = cn("base-class", false && "hidden-class", "visible-class");
    expect(result).toContain("base-class");
    expect(result).toContain("visible-class");
    expect(result).not.toContain("hidden-class");
  });

  test("should handle undefined and null", () => {
    const result = cn("base-class", undefined, null, "other-class");
    expect(result).toContain("base-class");
    expect(result).toContain("other-class");
  });

  test("should merge conflicting Tailwind classes", () => {
    // twMerge should handle conflicting classes by keeping the last one
    const result = cn("p-4", "p-8");
    expect(result).toBe("p-8");
  });

  test("should handle array of classes", () => {
    const result = cn(["text-red-500", "bg-blue-500"]);
    expect(result).toContain("text-red-500");
    expect(result).toContain("bg-blue-500");
  });

  test("should handle object with boolean values", () => {
    const result = cn({
      "text-red-500": true,
      "bg-blue-500": false,
      "border-gray-300": true,
    });
    expect(result).toContain("text-red-500");
    expect(result).not.toContain("bg-blue-500");
    expect(result).toContain("border-gray-300");
  });

  test("should handle empty input", () => {
    const result = cn();
    expect(result).toBe("");
  });

  test("should handle complex mixed inputs", () => {
    const isActive = true;
    const result = cn(
      "base-class",
      ["array-class-1", "array-class-2"],
      {
        "object-class-1": true,
        "object-class-2": false,
      },
      isActive && "conditional-class",
      "final-class"
    );

    expect(result).toContain("base-class");
    expect(result).toContain("array-class-1");
    expect(result).toContain("array-class-2");
    expect(result).toContain("object-class-1");
    expect(result).not.toContain("object-class-2");
    expect(result).toContain("conditional-class");
    expect(result).toContain("final-class");
  });

  test("should handle whitespace in class names", () => {
    const result = cn("  text-red-500  ", "  bg-blue-500  ");
    expect(result).toContain("text-red-500");
    expect(result).toContain("bg-blue-500");
  });
});
