<script lang="ts">
  import { HugeiconsIcon } from "@hugeicons/svelte";
  import { DatabaseExportIcon, DatabaseImportIcon } from "@hugeicons/core-free-icons";
  import { Button } from "$lib/components/ui/button/index.js";
  import { Label } from "$lib/components/ui/label/index.js";
  import * as RadioGroup from "$lib/components/ui/radio-group/index.js";
  import { toast } from "svelte-sonner";
  import { Export, Import } from "wailsjs/go/service/BackupService";

  let mode = $state<"export" | "import">("export");
  let processing = $state(false);

  async function handleAction() {
    processing = true;
    try {
      if (mode === "export") {
        await Export();
        toast.success("Backup exportado correctamente");
      } else {
        await Import();
        toast.success("Backup importado correctamente");
      }
    } catch (e) {
      toast.error(String(e));
    } finally {
      processing = false;
    }
  }
</script>

<div class="space-y-6">
  <h2 class="text-lg font-semibold">Backup</h2>

  <div class="max-w-md space-y-4">
    <RadioGroup.Root bind:value={mode} class="flex gap-4">
      <div class="flex items-center gap-2">
        <RadioGroup.Item value="export" id="export" />
        <Label for="export" class="cursor-pointer">Exportar</Label>
      </div>
      <div class="flex items-center gap-2">
        <RadioGroup.Item value="import" id="import" />
        <Label for="import" class="cursor-pointer">Importar</Label>
      </div>
    </RadioGroup.Root>

    {#if mode === "import"}
      <p class="text-sm text-muted-foreground">
        Seleccioná el archivo de base de datos (.db) de backup a importar. Esto reemplazará todos los datos actuales.
      </p>
    {/if}

    <Button onclick={handleAction} disabled={processing} class="w-full">
      <HugeiconsIcon icon={mode === "export" ? DatabaseExportIcon : DatabaseImportIcon} size={16} />
      {#if processing}
        {mode === "export" ? "Exportando..." : "Importando..."}
      {:else}
        {mode === "export" ? "Exportar Backup" : "Importar Backup"}
      {/if}
    </Button>
  </div>
</div>
