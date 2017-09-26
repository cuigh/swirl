///<reference path="../core/core.ts" />
namespace Swirl.Config {
    import Modal = Swirl.Core.Modal;
    import AjaxResult = Swirl.Core.AjaxResult;
    import Table = Swirl.Core.ListTable;

    export class ListPage {
        private table: Table;

        constructor() {
            this.table = new Table("#table-items");

            // bind events
            this.table.on("delete-config", this.deleteConfig.bind(this));
            $("#btn-delete").click(this.deleteConfigs.bind(this));
        }

        private deleteConfig(e: JQueryEventObject) {
            let $tr = $(e.target).closest("tr");
            let name = $tr.find("td:eq(1)").text().trim();
            let id = $tr.find(":checkbox:first").val();
            Modal.confirm(`Are you sure to remove config: <strong>${name}</strong>?`, "Delete config", (dlg, e) => {
                $ajax.post("delete", {ids: id}).trigger(e.target).encoder("form").json<AjaxResult>(r => {
                    $tr.remove();
                    dlg.close();
                })
            });
        }

        private deleteConfigs(e: JQueryEventObject) {
            let ids = this.table.selectedKeys();
            if (ids.length == 0) {
                Modal.alert("Please select one or more items.");
                return;
            }

            Modal.confirm(`Are you sure to remove ${ids.length} configs?`, "Delete configs", (dlg, e) => {
                $ajax.post("delete", {ids: ids.join(",")}).trigger(e.target).encoder("form").json<AjaxResult>(r => {
                    this.table.selectedRows().remove();
                    dlg.close();
                })
            });
        }
    }
}