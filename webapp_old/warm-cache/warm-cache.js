var kvik = "http://kvik.cs.uit.no"; 
var names = [ 
"hsa01210",
"hsa02010",
"hsa04152",
"hsa05221",
"hsa04520",
"hsa04920",
"hsa04261",
"hsa05143",
"hsa00250",
"hsa05034",
"hsa04960",
"hsa05330",
"hsa05010",
"hsa00520",
"hsa00970",
"hsa05146",
"hsa05031",
"hsa05014",
"hsa04612",
"hsa04210",
"hsa00590",
"hsa00330",
"hsa05412",
"hsa00053",
"hsa05310",
"hsa05320",
"hsa04360",
"hsa04662",
"hsa05100",
"hsa05217",
"hsa03022",
"hsa03410",
"hsa04976",
"hsa01230",
"hsa01040",
"hsa00780",
"hsa05219",
"hsa00650",
"hsa00524",
"hsa00232",
"hsa04020",
"hsa04973",
"hsa01200",
"hsa04260",
"hsa04514",
"hsa04110",
"hsa05230",
"hsa05142",
"hsa05204",
"hsa04062",
"hsa05231",
"hsa04725",
"hsa05220",
"hsa04713",
"hsa04710",
"hsa00020",
"hsa05030",
"hsa04966",
"hsa05210",
"hsa04610",
"hsa00460",
"hsa00270",
"hsa04060",
"hsa04623",
"hsa00472",
"hsa00471",
"hsa03030",
"hsa01220",
"hsa05414",
"hsa04728",
"hsa04320",
"hsa00982",
"hsa00983",
"hsa04512",
"hsa04961",
"hsa04144",
"hsa05213",
"hsa05120",
"hsa05169",
"hsa04012",
"hsa04915",
"hsa00565",
"hsa03460",
"hsa04975",
"hsa00061",
"hsa00071",
"hsa00062",
"hsa01212",
"hsa04664",
"hsa04666",
"hsa04510",
"hsa00790",
"hsa04068",
"hsa00051",
"hsa04727",
"hsa00052",
"hsa04540",
"hsa04971",
"hsa05214",
"hsa04724",
"hsa00480",
"hsa00561",
"hsa00564",
"hsa00260",
"hsa00010",
"hsa00532",
"hsa00534",
"hsa00533",
"hsa00531",
"hsa00604",
"hsa00603",
"hsa00601",
"hsa00563",
"hsa00630",
"hsa04912",
"hsa05332",
"hsa04066",
"hsa05166",
"hsa04340",
"hsa04640",
"hsa05161",
"hsa05160",
"hsa05168",
"hsa04390",
"hsa00340",
"hsa03440",
"hsa05016",
"hsa05410",
"hsa05321",
"hsa04750",
"hsa05164",
"hsa00562",
"hsa04911",
"hsa04910",
"hsa04672",
"hsa04630",
"hsa05134",
"hsa05140",
"hsa04670",
"hsa00591",
"hsa00785",
"hsa04730",
"hsa04720",
"hsa00300",
"hsa00310",
"hsa04142",
"hsa04010",
"hsa05144",
"hsa04950",
"hsa05162",
"hsa04916",
"hsa05218",
"hsa01100",
"hsa00980",
"hsa05206",
"hsa04978",
"hsa03430",
"hsa05032",
"hsa00512",
"hsa00510",
"hsa04064",
"hsa04621",
"hsa04650",
"hsa04080",
"hsa04722",
"hsa00760",
"hsa05033",
"hsa00910",
"hsa04932",
"hsa03450",
"hsa05223",
"hsa04330",
"hsa03420",
"hsa04740",
"hsa00670",
"hsa04114",
"hsa04380",
"hsa00511",
"hsa00514",
"hsa04913",
"hsa00190",
"hsa04921",
"hsa04151",
"hsa03320",
"hsa05212",
"hsa04972",
"hsa00770",
"hsa05012",
"hsa05130",
"hsa05200",
"hsa00040",
"hsa00030",
"hsa04146",
"hsa05133",
"hsa04145",
"hsa00360",
"hsa00400",
"hsa04070",
"hsa04744",
"hsa04611",
"hsa00860",
"hsa00120",
"hsa05340",
"hsa05020",
"hsa04914",
"hsa04917",
"hsa00640",
"hsa05215",
"hsa03050",
"hsa04974",
"hsa03060",
"hsa04141",
"hsa05205",
"hsa04964",
"hsa00230",
"hsa00240",
"hsa00620",
"hsa04622",
"hsa03018",
"hsa03020",
"hsa03013",
"hsa04015",
"hsa04014",
"hsa04810",
"hsa04140",
"hsa05211",
"hsa04614",
"hsa00830",
"hsa04723",
"hsa05323",
"hsa00740",
"hsa03010",
"hsa03008",
"hsa04130",
"hsa04970",
"hsa05132",
"hsa00450",
"hsa04726",
"hsa05131",
"hsa04550",
"hsa05222",
"hsa00600",
"hsa03040",
"hsa05150",
"hsa00500",
"hsa00100",
"hsa00140",
"hsa00920",
"hsa04122",
"hsa04721",
"hsa00072",
"hsa05322",
"hsa04660",
"hsa04350",
"hsa04668",
"hsa04742",
"hsa00430",
"hsa00900",
"hsa00730",
"hsa05216",
"hsa04919",
"hsa04918",
"hsa04530",
"hsa04620",
"hsa05145",
"hsa05202",
"hsa00380",
"hsa05152",
"hsa04940",
"hsa04930",
"hsa00350",
"hsa00130",
"hsa04120",
"hsa04370",
"hsa00290",
"hsa00280",
"hsa04270",
"hsa04962",
"hsa05110",
"hsa05203",
"hsa05416",
"hsa00750",
"hsa04977",
"hsa04310",
"hsa00592",
"hsa00410",
"hsa04024",
"hsa04022",
"hsa03015",
"hsa04150",
"hsa04115"
]
module.exports = {
    'Pathway browser is up': function(test) { 
        test.open(kvik + "/browser").done(); 
    }, 
    'All the pathways': function(test){ 
        test.open(kvik + "/browser/pathwaySelect="+"hsa01210").done(); 
        test.open(kvik + "/browser/pathwaySelect="+"hsa02010").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04152").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05221").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04520").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04920").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04261").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05143").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00250").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05034").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04960").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05330").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05010").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00520").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00970").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05146").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05031").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05014").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04612").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04210").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00590").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00330").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05412").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00053").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05310").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05320").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04360").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04662").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05100").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05217").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa03022").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa03410").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04976").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa01230").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa01040").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00780").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05219").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00650").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00524").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00232").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04020").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04973").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa01200").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04260").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04514").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04110").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05230").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05142").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05204").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04062").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05231").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04725").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05220").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04713").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04710").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00020").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05030").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04966").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05210").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04610").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00460").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00270").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04060").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04623").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00472").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00471").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa03030").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa01220").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05414").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04728").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04320").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00982").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00983").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04512").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04961").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04144").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05213").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05120").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05169").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04012").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04915").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00565").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa03460").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04975").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00061").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00071").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00062").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa01212").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04664").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04666").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04510").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00790").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04068").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00051").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04727").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00052").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04540").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04971").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05214").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04724").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00480").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00561").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00564").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00260").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00010").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00532").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00534").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00533").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00531").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00604").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00603").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00601").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00563").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00630").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04912").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05332").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04066").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05166").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04340").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04640").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05161").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05160").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05168").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04390").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00340").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa03440").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05016").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05410").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05321").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04750").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05164").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00562").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04911").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04910").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04672").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04630").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05134").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05140").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04670").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00591").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00785").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04730").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04720").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00300").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00310").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04142").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04010").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05144").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04950").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05162").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04916").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05218").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa01100").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00980").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05206").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04978").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa03430").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05032").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00512").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00510").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04064").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04621").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04650").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04080").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04722").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00760").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05033").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00910").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04932").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa03450").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05223").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04330").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa03420").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04740").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00670").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04114").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04380").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00511").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00514").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04913").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00190").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04921").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04151").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa03320").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05212").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04972").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00770").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05012").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05130").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05200").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00040").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00030").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04146").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05133").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04145").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00360").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00400").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04070").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04744").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04611").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00860").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00120").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05340").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05020").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04914").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04917").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00640").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05215").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa03050").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04974").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa03060").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04141").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05205").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04964").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00230").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00240").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00620").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04622").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa03018").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa03020").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa03013").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04015").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04014").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04810").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04140").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05211").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04614").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00830").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04723").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05323").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00740").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa03010").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa03008").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04130").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04970").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05132").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00450").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04726").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05131").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04550").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05222").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00600").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa03040").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05150").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00500").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00100").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00140").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00920").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04122").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04721").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00072").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05322").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04660").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04350").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04668").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04742").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00430").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00900").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00730").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05216").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04919").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04918").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04530").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04620").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05145").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05202").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00380").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05152").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04940").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04930").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00350").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00130").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04120").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04370").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00290").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00280").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04270").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04962").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05110").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05203").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa05416").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00750").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04977").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04310").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00592").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa00410").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04024").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04022").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa03015").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04150").done();
        test.open(kvik + "/browser/pathwaySelect="+"hsa04115").done();
    }




};
