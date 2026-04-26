<script lang="ts">
  import { HugeiconsIcon } from "@hugeicons/svelte";
  import { PencilEdit01Icon, Delete01Icon, PlusSignIcon, Download01Icon } from "@hugeicons/core-free-icons";
  import { Button } from "$lib/components/ui/button/index.js";
  import { Input } from "$lib/components/ui/input/index.js";
  import * as Table from "$lib/components/ui/table/index.js";
  import { toast } from "svelte-sonner";
  import { ExportZones } from "wailsjs/go/service/ExportService";
  import { getZones, getZonesLoading, fetchZones } from "$lib/stores/zones.svelte.js";
  import ZoneFormDialog from "./ZoneFormDialog.svelte";
  import ZoneDeleteDialog from "./ZoneDeleteDialog.svelte";
  import type { model } from "wailsjs/go/models";

  async function exportZones() {
    try {
      await ExportZones();
      toast.success("Zonas exportadas correctamente");
    } catch (e) {
      toast.error(String(e));
    }
  }

  let formOpen = $state(false);
  let deleteOpen = $state(false);
  let selectedZone = $state<model.Zone | undefined>(undefined);
  let searchQuery = $state("");

  let zones = $derived(getZones());
  let loading = $derived(getZonesLoading());
  let filteredZones = $derived(
    zones.filter((z) => z.name.toLowerCase().includes(searchQuery.toLowerCase()))
  );

  $effect(() => {
    fetchZones();
  });

  function openCreate() {
    selectedZone = undefined;
    formOpen = true;
  }

  function openEdit(zone: model.Zone) {
    selectedZone = zone;
    formOpen = true;
  }

  function openDelete(zone: model.Zone) {
    selectedZone = zone;
    deleteOpen = true;
  }
</script>

<div class="space-y-4">
  <div class="flex items-center justify-between">
    <h2 class="text-lg font-semibold">Zonas</h2>
    <div class="flex gap-2">
      <Button variant="outline" size="default" onclick={exportZones}>
        <HugeiconsIcon icon={Download01Icon} size={14} />
        Exportar
      </Button>
      <Button size="default" onclick={openCreate}>
        <HugeiconsIcon icon={PlusSignIcon} size={14} />
        Nueva Zona
      </Button>
    </div>
  </div>

  <Input placeholder="Buscar por nombre..." bind:value={searchQuery} class="max-w-xs" />

  {#if loading && zones.length === 0}
    <p class="text-sm text-muted-foreground">Cargando...</p>
  {:else if filteredZones.length === 0}
    <p class="text-sm text-muted-foreground">No hay zonas registradas.</p>
  {:else}
    <Table.Root>
      <Table.Header>
        <Table.Row>
          <Table.Head>Nombre</Table.Head>
          <Table.Head class="w-32 text-right">Acciones</Table.Head>
        </Table.Row>
      </Table.Header>
      <Table.Body>
        {#each filteredZones as zone (zone.id)}
          <Table.Row>
            <Table.Cell><span class="text-pretty">{zone.name}</span></Table.Cell>
            <Table.Cell class="text-right">
              <div class="flex justify-end gap-1">
                <Button title="Editar" variant="ghost" size="icon" onclick={() => openEdit(zone)}>
                  <HugeiconsIcon icon={PencilEdit01Icon} size={14} />
                </Button>
                <Button title="Eliminar" variant="ghost" size="icon" onclick={() => openDelete(zone)}>
                  <HugeiconsIcon icon={Delete01Icon} size={14} />
                </Button>
              </div>
            </Table.Cell>
          </Table.Row>
        {/each}
      </Table.Body>
    </Table.Root>
  {/if}
</div>

<ZoneFormDialog bind:open={formOpen} zone={selectedZone} />

{#if selectedZone}
  <ZoneDeleteDialog bind:open={deleteOpen} zone={selectedZone} />
{/if}
