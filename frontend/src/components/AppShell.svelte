<script lang="ts">
  import { HugeiconsIcon } from "@hugeicons/svelte";
  import { MapsLocation01Icon, UserMultipleIcon, Package01Icon, DeliveryTruck01Icon, Analytics01Icon, ArrowMoveUpLeftIcon } from "@hugeicons/core-free-icons";
  import { cn } from "$lib/utils.js";
  import ZonesPage from "./zones/ZonesPage.svelte";
  import ClientsPage from "./clients/ClientsPage.svelte";
  import ProductsPage from "./products/ProductsPage.svelte";
  import DeliveryPage from "./delivery/DeliveryPage.svelte";
  import ZoneReportsPage from "./zones/ZoneReportsPage.svelte";
  import ZoneReassignPage from "./zones/ZoneReassignPage.svelte";

  type Page = "zones" | "clients" | "products" | "delivery" | "reports" | "reassign";
  let currentPage = $state<Page>("zones");

  const navItems: { page: Page; label: string; icon: typeof MapsLocation01Icon }[] = [
    { page: "zones", label: "Zonas", icon: MapsLocation01Icon },
    { page: "clients", label: "Clientes", icon: UserMultipleIcon },
    { page: "products", label: "Productos", icon: Package01Icon },
    { page: "delivery", label: "Entregas", icon: DeliveryTruck01Icon },
    { page: "reports", label: "Reportes", icon: Analytics01Icon },
    { page: "reassign", label: "Reasignar", icon: ArrowMoveUpLeftIcon },
  ];
</script>

<div class="flex h-screen">
  <nav class="print-hidden flex w-48 flex-col gap-1 border-r bg-muted/30 p-3">
    <h1 class="mb-4 px-2 text-sm font-semibold tracking-tight">SoditApp</h1>
    {#each navItems as item}
      <button
        class={cn(
          "flex items-center gap-2 rounded-md px-2 py-1.5 text-sm transition-colors",
          currentPage === item.page
            ? "bg-primary text-primary-foreground"
            : "text-muted-foreground hover:bg-muted hover:text-foreground"
        )}
        onclick={() => (currentPage = item.page)}
      >
        <HugeiconsIcon icon={item.icon} size={16} />
        {item.label}
      </button>
    {/each}
  </nav>

  <main class="flex-1 overflow-auto p-6">
    {#if currentPage === "zones"}
      <ZonesPage />
    {:else if currentPage === "clients"}
      <ClientsPage />
    {:else if currentPage === "products"}
      <ProductsPage />
    {:else if currentPage === "delivery"}
      <DeliveryPage />
    {:else if currentPage === "reports"}
      <ZoneReportsPage />
    {:else if currentPage === "reassign"}
      <ZoneReassignPage />
    {/if}
  </main>
</div>
