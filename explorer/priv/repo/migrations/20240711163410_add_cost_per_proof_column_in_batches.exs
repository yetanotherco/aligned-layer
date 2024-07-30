defmodule Explorer.Repo.Migrations.AddCostPerProofColumnInBatches do
  use Ecto.Migration

  def change do
    alter table("batches") do
      add :cost_per_proof, :bigint
    end
  end
end
