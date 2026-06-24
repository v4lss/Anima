import { api } from "./api";

// Mock fetch for testing
global.fetch = jest.fn();

describe("API service", () => {
  beforeEach(() => {
    (global.fetch as jest.Mock).mockClear();
  });

  test("get makes GET request with correct URL", async () => {
    (global.fetch as jest.Mock).mockResolvedValue({
      ok: true,
      json: async () => ({ data: [] }),
    });

    await api.get("/api/monitors");

    expect(global.fetch).toHaveBeenCalledWith(
      expect.stringContaining("/api/monitors"),
      expect.objectContaining({
        method: "GET",
        headers: expect.objectContaining({
          "Content-Type": "application/json",
        }),
      })
    );
  });

  test("post makes POST request with body", async () => {
    (global.fetch as jest.Mock).mockResolvedValue({
      ok: true,
      json: async () => ({ data: {} }),
    });

    await api.post("/api/auth/login", { email: "test@example.com", password: "pass" });

    expect(global.fetch).toHaveBeenCalledWith(
      expect.stringContaining("/api/auth/login"),
      expect.objectContaining({
        method: "POST",
        body: expect.stringContaining("test@example.com"),
      })
    );
  });

  test("delete makes DELETE request", async () => {
    (global.fetch as jest.Mock).mockResolvedValue({
      ok: true,
      json: async () => ({ data: {} }),
    });

    await api.delete("/api/monitors/123");

    expect(global.fetch).toHaveBeenCalledWith(
      expect.stringContaining("/api/monitors/123"),
      expect.objectContaining({
        method: "DELETE",
      })
    );
  });

  test("handles error response", async () => {
    (global.fetch as jest.Mock).mockResolvedValue({
      ok: false,
      json: async () => ({ error: "Unauthorized" }),
    });

    await expect(api.get("/api/monitors")).rejects.toThrow("Unauthorized");
  });
});
