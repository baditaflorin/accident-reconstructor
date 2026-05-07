import { expect, test } from "@playwright/test";

test("loads the Pages app and renders a sample 3D reconstruction", async ({
  page,
}) => {
  await page.goto("./");

  await expect(
    page.getByRole("heading", { name: "Accident Reconstructor" }),
  ).toBeVisible();
  await expect(
    page.getByRole("link", { name: "Star on GitHub" }),
  ).toHaveAttribute(
    "href",
    "https://github.com/baditaflorin/accident-reconstructor",
  );

  await page.getByRole("button", { name: "Load Sample" }).click();
  await expect(page.getByText("27.9 km/h")).toBeVisible();
  await expect(page.locator("canvas")).toBeVisible();
});
