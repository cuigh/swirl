///<reference path="../core/core.ts" />
namespace Swirl.Volume {
    import OptionTable = Swirl.Core.OptionTable;

    export class NewPage {
        constructor() {
            new OptionTable("#table-options");
            new OptionTable("#table-labels");

            $("#drivers").find(":radio[name=driver]").change(e => {
                $("#txt-custom-driver").prop("disabled", $(e.target).val() != "other");
            });
        }
    }
}