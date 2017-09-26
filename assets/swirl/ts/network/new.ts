///<reference path="../core/core.ts" />
namespace Swirl.Network {
    import OptionTable = Swirl.Core.OptionTable;

    export class NewPage {
        constructor() {
            new OptionTable("#table-options");
            new OptionTable("#table-labels");

            $("#drivers :radio[name=driver]").change(e => {
                $("#txt-custom-driver").prop("disabled", $(e.target).val() != "other");
            });
            $("#ipv6_enabled").change(e => {
                let enabled = $(e.target).prop("checked");
                $("#ipv6_subnet,#ipv6_gateway").prop("disabled", !enabled);
            });
        }
    }
}