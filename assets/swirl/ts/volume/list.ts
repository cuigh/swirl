///<reference path="../core/core.ts" />
namespace Swirl.Volume {
    import Modal = Swirl.Core.Modal;
    import AjaxResult = Swirl.Core.AjaxResult;
    import Table = Swirl.Core.ListTable;

    export class ListPage {
        private table: Table;

        constructor() {
            this.table = new Table("#table-items");

            // bind events
            this.table.on("delete-volume", this.deleteVolume.bind(this));
            $("#btn-delete").click(this.deleteVolumes.bind(this));
            $("#btn-prune").click(this.pruneVolumes.bind(this));
        }

        private deleteVolume(e: JQueryEventObject) {
            let $tr = $(e.target).closest("tr");
            let name = $tr.find("td:eq(1)").text().trim();
            Modal.confirm(`Are you sure to remove volume: <strong>${name}</strong>?`, "Delete volume", (dlg, e) => {
                $ajax.post("delete", { names: name }).trigger(e.target).encoder("form").json<AjaxResult>(r => {
                    $tr.remove();
                    dlg.close();
                })
            });
        }

        private deleteVolumes(e: JQueryEventObject) {
            let names = this.table.selectedKeys();
            if (names.length == 0) {
                Modal.alert("Please select one or more items.")
                return;
            }

            Modal.confirm(`Are you sure to remove ${names.length} volumes?`, "Delete volumes", (dlg, e) => {
                $ajax.post("delete", { names: names.join(",") }).trigger(e.target).encoder("form").json<AjaxResult>(r => {
                    this.table.selectedRows().remove();
                    dlg.close();
                })
            });
        }

        private pruneVolumes(e: JQueryEventObject) {
            Modal.confirm(`Are you sure to remove all unused volumes?`, "Prune volumes", (dlg, e) => {
                $ajax.post("prune").trigger(e.target).json<AjaxResult>(r => {
                    location.reload();
                })
            });
        }
    }
}