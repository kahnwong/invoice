(function() {
    var invoicer = angular.module('invoicer', ["xeditable"]);

    invoicer.run(function(editableOptions) {
      editableOptions.theme = 'bs3';
      });

    invoicer.controller('invController', function($scope){
        $scope.invoice = work;

        $scope.getSubTotal = function(){
            var subtotal=0;
            for(var i=0; i < $scope.invoice.line_items.length; i++) {
                var line = $scope.invoice.line_items[i];
                subtotal += (line.hours * line.unit_price);
            }
            return subtotal;
        }

        $scope.getTax = function(subtotal){
            var tax = 0;
            tax = subtotal * 0.07;
            return tax;
        }

        $scope.addItem = function() {
            $scope.invoice.line_items.push(this.invoice.temp);
            $scope.invoice.temp = {};
        }

       $scope.contentLoaded = true;

    });

    var work =
    {
        bill_type: "ใบกำกับภาษี",
        tax_id: "00000000000",
        number: "1234",
        send_date: "15 กันยายน, 2560",
        due_date: "25 กันยายน, 2560",
        from_company: "ชัยกลกิจ",
        from_address: "730/10 ถ. จันท์ บางโคล่",
        from_city: "ยานนาวา",
        from_state: "กรุงเทพฯ",
        from_zip: "10120",
        from_phone: "081-869-8851",
        from_email: "chaikolkit@gmail.com",
        to_company: "บริษัทผลิตภัณฑ์ตราเพชร จำกัด มหาชน",
        to_address: "69-70 หมู่ 1 ถ. มิตรภาพ ตลิ่งชัน",
        to_city: "เมือง",
        to_state: "สระบุรี",
        to_zip: "18000",
        to_phone: "(555) 111-2222",
        to_email: "their.email@address.com",
        line_items: [
            {
                title: "Item 1",
                desc: "ขนาด 2x10",
                hours: "2",
                unit_price: "200"
            },
            {
                title: "Item 2",
                desc: "ขนาด 2x5",
                hours: "2",
                unit_price: "50"
            },
        ]
    };

})();
