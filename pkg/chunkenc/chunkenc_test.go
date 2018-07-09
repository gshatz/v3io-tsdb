/*
Copyright 2018 Iguazio Systems Ltd.

Licensed under the Apache License, Version 2.0 (the "License") with
an addition restriction as set forth herein. You may not use this
file except in compliance with the License. You may obtain a copy of
the License at http://www.apache.org/licenses/LICENSE-2.0.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
implied. See the License for the specific language governing
permissions and limitations under the License.

In addition, you may not use the software for any purposes that are
illegal under applicable law, and the grant of the foregoing license
under the Apache 2.0 license is conditioned upon your compliance with
such restriction.
*/

package chunkenc

import (
	"fmt"
	"testing"

	"encoding/base64"
	"math/rand"
	"time"
)

const basetime = 1524690488000

type sample struct {
	t int64
	v float64
}

// [132 180 199 187 191 88 63 240 - 0 0 0 0 0 0 154 8 - 194 95 255 108 7 126 113 172 - 46 18 195 104 59 202 237 129 - 119 243 146]

func TestXor(tst *testing.T) {

	samples := GenSamples(1000, 5, 1000, 100)
	//samples := RealSample(1000)
	byteArray := []byte{}

	ch := NewXORChunk()
	appender, err := ch.Appender()
	if err != nil {
		tst.Fatal(err)
	}

	for i, s := range samples {
		fmt.Println("t,v: ", s.t, s.v)
		appender.Append(s.t, s.v)
		b := ch.Bytes()
		fmt.Println(b, len(b))
		byteArray = append(byteArray, b...)
		ch.Clear()
		if i == 4 {
			fmt.Println("restarted appender")
			ch = NewXORChunk()
			appender, err = ch.Appender()
			if err != nil {
				tst.Fatal(err)
			}

		}
	}

	fmt.Println("Samples:", len(samples), "byteArray:", byteArray, len(byteArray))

	ch2, err := FromData(EncXOR, byteArray, 0)
	if err != nil {
		tst.Fatal(err)
	}

	iter := ch2.Iterator()
	i := 0
	for iter.Next() {

		if iter.Err() != nil {
			tst.Fatal(iter.Err())
		}

		t, v := iter.At()
		isMatch := t == samples[i].t && v == samples[i].v
		fmt.Println("t, v, match: ", t, v, isMatch)
		if !isMatch {
			tst.Fatalf("iterator t or v doesnt match appended index %d len %d", i, len(samples))
		}
		i++
	}
	fmt.Println()

	if i != len(samples) {
		tst.Fatalf("number of iterator samples (%d) != num of appended (%d)", i, len(samples))
	}

}

func TestBstream(t *testing.T) {
	src := &bstream{count: 8, stream: []byte{0x55, 0x44, 0x33}}

	bs := newBWriter(8)
	byt, _ := src.readByte()
	bs.writeByte(byt)
	fmt.Println(bs.count, bs.stream)
	for i := 1; i < 18; i++ {
		bit, _ := src.readBit()
		fmt.Println(bs.count, bs.stream, bit)
		bs.writeBit(bit)
	}

	fmt.Println("Reading:")
	bs2 := &bstream{count: 8, stream: bs.stream}
	fmt.Println(bs2.count, bs2.stream)
	for i := 1; i < 18; i++ {
		bit, _ := bs2.readBit()
		fmt.Println(bs2.count, bs2.stream, bit)
	}

}

func DecodeTest(blob string) error {

	//blob := "+AFjT7+iCEBLgAAAAAAA+AFjT8A6YEBQgAAAAAAA"

	data, err := base64.StdEncoding.DecodeString(blob)
	if err != nil {
		return err
	}
	fmt.Println(data)

	chunk, err := FromData(EncXOR, data, 0)
	if err != nil {
		return err
	}

	iter := chunk.Iterator()
	i := 0
	for iter.Next() {

		if iter.Err() != nil {
			return iter.Err()
		}

		t, v := iter.At()
		tstr := time.Unix(int64(t/1000), 0).Format(time.RFC3339)
		fmt.Printf("unix=%d, t=%s, v=%.4f \n", t, tstr, v)
		i++
	}

	return nil
}

func GenSamples(num, interval int, start, step float64) []sample {
	samples := []sample{}
	curTime := int64(basetime)
	v := start

	for i := 0; i <= num; i++ {
		curTime += int64(interval * 1000)
		t := curTime + int64(rand.Intn(100)) - 50
		v += float64(rand.Intn(100)-50) / 100 * step
		//fmt.Printf("t-%d,v%.2f ", t, v)
		samples = append(samples, sample{t: t, v: v})
	}

	return samples
}

var timeList = []int64{1360281600000, 1360540800000, 1360627200000, 1360713600000, 1360800000000, 1360886400000, 1361232000000, 1361318400000, 1361404800000, 1361491200000, 1361750400000, 1361836800000, 1361923200000, 1362009600000, 1362096000000, 1362355200000, 1362441600000, 1362528000000, 1362614400000, 1362700800000, 1362960000000, 1363046400000, 1363132800000, 1363219200000, 1363305600000, 1363564800000, 1363651200000, 1363737600000, 1363824000000, 1363910400000, 1364169600000, 1364256000000, 1364342400000, 1364428800000, 1364774400000, 1364860800000, 1364947200000, 1365033600000, 1365120000000, 1365379200000, 1365465600000, 1365552000000, 1365638400000, 1365724800000, 1365984000000, 1366070400000, 1366156800000, 1366243200000, 1366329600000, 1366588800000, 1366675200000, 1366761600000, 1366848000000, 1366934400000, 1367193600000, 1367280000000, 1367366400000, 1367452800000, 1367539200000, 1367798400000, 1367884800000, 1367971200000, 1368057600000, 1368144000000, 1368403200000, 1368489600000, 1368576000000, 1368662400000, 1368748800000, 1369008000000, 1369094400000, 1369180800000, 1369267200000, 1369353600000, 1369699200000, 1369785600000, 1369872000000, 1369958400000, 1370217600000, 1370304000000, 1370390400000, 1370476800000, 1370563200000, 1370822400000, 1370908800000, 1370995200000, 1371081600000, 1371168000000, 1371427200000, 1371513600000, 1371600000000, 1371686400000, 1371772800000, 1372032000000, 1372118400000, 1372204800000, 1372291200000, 1372377600000, 1372636800000, 1372723200000, 1372809600000, 1372982400000, 1373241600000, 1373328000000, 1373414400000, 1373500800000, 1373587200000, 1373846400000, 1373932800000, 1374019200000, 1374105600000, 1374192000000, 1374451200000, 1374537600000, 1374624000000, 1374710400000, 1374796800000, 1375056000000, 1375142400000, 1375228800000, 1375315200000, 1375401600000, 1375660800000, 1375747200000, 1375833600000, 1375920000000, 1376006400000, 1376265600000, 1376352000000, 1376438400000, 1376524800000, 1376611200000, 1376870400000, 1376956800000, 1377043200000, 1377129600000, 1377216000000, 1377475200000, 1377561600000, 1377648000000, 1377734400000, 1377820800000, 1378166400000, 1378252800000, 1378339200000, 1378425600000, 1378684800000, 1378771200000, 1378857600000, 1378944000000, 1379030400000, 1379289600000, 1379376000000, 1379462400000, 1379548800000, 1379635200000, 1379894400000, 1379980800000, 1380067200000, 1380153600000, 1380240000000, 1380499200000, 1380585600000, 1380672000000, 1380758400000, 1380844800000, 1381104000000, 1381190400000, 1381276800000, 1381363200000, 1381449600000, 1381708800000, 1381795200000, 1381881600000, 1381968000000, 1382054400000, 1382313600000, 1382400000000, 1382486400000, 1382572800000, 1382659200000, 1382918400000, 1383004800000, 1383091200000, 1383177600000, 1383264000000, 1383523200000, 1383609600000, 1383696000000, 1383782400000, 1383868800000, 1384128000000, 1384214400000, 1384300800000, 1384387200000, 1384473600000, 1384732800000, 1384819200000, 1384905600000, 1384992000000, 1385078400000, 1385337600000, 1385424000000, 1385510400000, 1385683200000, 1385942400000, 1386028800000, 1386115200000, 1386201600000, 1386288000000, 1386547200000, 1386633600000, 1386720000000, 1386806400000, 1386892800000, 1387152000000, 1387238400000, 1387324800000, 1387411200000, 1387497600000, 1387756800000, 1387843200000, 1388016000000, 1388102400000, 1388361600000, 1388448000000, 1388620800000, 1388707200000, 1388966400000, 1389052800000, 1389139200000, 1389225600000, 1389312000000, 1389571200000, 1389657600000, 1389744000000, 1389830400000, 1389916800000, 1390262400000, 1390348800000, 1390435200000, 1390521600000, 1390780800000, 1390867200000, 1390953600000, 1391040000000, 1391126400000, 1391385600000, 1391472000000, 1391558400000, 1391644800000, 1391731200000, 1391990400000, 1392076800000, 1392163200000, 1392249600000, 1392336000000, 1392681600000, 1392768000000, 1392854400000, 1392940800000, 1393200000000, 1393286400000, 1393372800000, 1393459200000, 1393545600000, 1393804800000, 1393891200000, 1393977600000, 1394064000000, 1394150400000, 1394409600000, 1394496000000, 1394582400000, 1394668800000, 1394755200000, 1395014400000, 1395100800000, 1395187200000, 1395273600000, 1395360000000, 1395619200000, 1395705600000, 1395792000000, 1395878400000, 1395964800000, 1396224000000, 1396310400000, 1396396800000, 1396483200000, 1396569600000, 1396828800000, 1396915200000, 1397001600000, 1397088000000, 1397174400000, 1397433600000, 1397520000000, 1397606400000, 1397692800000, 1398038400000, 1398124800000, 1398211200000, 1398297600000, 1398384000000, 1398643200000, 1398729600000, 1398816000000, 1398902400000, 1398988800000, 1399248000000, 1399334400000, 1399420800000, 1399507200000, 1399593600000, 1399852800000, 1399939200000, 1400025600000, 1400112000000, 1400198400000, 1400457600000, 1400544000000, 1400630400000, 1400716800000, 1400803200000, 1401148800000, 1401235200000, 1401321600000, 1401408000000, 1401667200000, 1401753600000, 1401840000000, 1401926400000, 1402012800000, 1402272000000, 1402358400000, 1402444800000, 1402531200000, 1402617600000, 1402876800000, 1402963200000, 1403049600000, 1403136000000, 1403222400000, 1403481600000, 1403568000000, 1403654400000, 1403740800000, 1403827200000, 1404086400000, 1404172800000, 1404259200000, 1404345600000, 1404691200000, 1404777600000, 1404864000000, 1404950400000, 1405036800000, 1405296000000, 1405382400000, 1405468800000, 1405555200000, 1405641600000, 1405900800000, 1405987200000, 1406073600000, 1406160000000, 1406246400000, 1406505600000, 1406592000000, 1406678400000, 1406764800000, 1406851200000, 1407110400000, 1407196800000, 1407283200000, 1407369600000, 1407456000000, 1407715200000, 1407801600000, 1407888000000, 1407974400000, 1408060800000, 1408320000000, 1408406400000, 1408492800000, 1408579200000, 1408665600000, 1408924800000, 1409011200000, 1409097600000, 1409184000000, 1409270400000, 1409616000000, 1409702400000, 1409788800000, 1409875200000, 1410134400000, 1410220800000, 1410307200000, 1410393600000, 1410480000000, 1410739200000, 1410825600000, 1410912000000, 1410998400000, 1411084800000, 1411344000000, 1411430400000, 1411516800000, 1411603200000, 1411689600000, 1411948800000, 1412035200000, 1412121600000, 1412208000000, 1412294400000, 1412553600000, 1412640000000, 1412726400000, 1412812800000, 1412899200000, 1413158400000, 1413244800000, 1413331200000, 1413417600000, 1413504000000, 1413763200000, 1413849600000, 1413936000000, 1414022400000, 1414108800000, 1414368000000, 1414454400000, 1414540800000, 1414627200000, 1414713600000, 1414972800000, 1415059200000, 1415145600000, 1415232000000, 1415318400000, 1415577600000, 1415664000000, 1415750400000, 1415836800000, 1415923200000, 1416182400000, 1416268800000, 1416355200000, 1416441600000, 1416528000000, 1416787200000, 1416873600000, 1416960000000, 1417132800000, 1417392000000, 1417478400000, 1417564800000, 1417651200000, 1417737600000, 1417996800000, 1418083200000, 1418169600000, 1418256000000, 1418342400000, 1418601600000, 1418688000000, 1418774400000, 1418860800000, 1418947200000, 1419206400000, 1419292800000, 1419379200000, 1419552000000, 1419811200000, 1419897600000, 1419984000000, 1420156800000, 1420416000000, 1420502400000, 1420588800000, 1420675200000, 1420761600000, 1421020800000, 1421107200000, 1421193600000, 1421280000000, 1421366400000, 1421712000000, 1421798400000, 1421884800000, 1421971200000, 1422230400000, 1422316800000, 1422403200000, 1422489600000, 1422576000000, 1422835200000, 1422921600000, 1423008000000, 1423094400000, 1423180800000, 1423440000000, 1423526400000, 1423612800000, 1423699200000, 1423785600000, 1424131200000, 1424217600000, 1424304000000, 1424390400000, 1424649600000, 1424736000000, 1424822400000, 1424908800000, 1424995200000, 1425254400000, 1425340800000, 1425427200000, 1425513600000, 1425600000000, 1425859200000, 1425945600000, 1426032000000, 1426118400000, 1426204800000, 1426464000000, 1426550400000, 1426636800000, 1426723200000, 1426809600000, 1427068800000, 1427155200000, 1427241600000, 1427328000000, 1427414400000, 1427673600000, 1427760000000, 1427846400000, 1427932800000, 1428278400000, 1428364800000, 1428451200000, 1428537600000, 1428624000000, 1428883200000, 1428969600000, 1429056000000, 1429142400000, 1429228800000, 1429488000000, 1429574400000, 1429660800000, 1429747200000, 1429833600000, 1430092800000, 1430179200000, 1430265600000, 1430352000000, 1430438400000, 1430697600000, 1430784000000, 1430870400000, 1430956800000, 1431043200000, 1431302400000, 1431388800000, 1431475200000, 1431561600000, 1431648000000, 1431907200000, 1431993600000, 1432080000000, 1432166400000, 1432252800000, 1432598400000, 1432684800000, 1432771200000, 1432857600000, 1433116800000, 1433203200000, 1433289600000, 1433376000000, 1433462400000, 1433721600000, 1433808000000, 1433894400000, 1433980800000, 1434067200000, 1434326400000, 1434412800000, 1434499200000, 1434585600000, 1434672000000, 1434931200000, 1435017600000, 1435104000000, 1435190400000, 1435276800000, 1435536000000, 1435622400000, 1435708800000, 1435795200000, 1436140800000, 1436227200000, 1436313600000, 1436400000000, 1436486400000, 1436745600000, 1436832000000, 1436918400000, 1437004800000, 1437091200000, 1437350400000, 1437436800000, 1437523200000, 1437609600000, 1437696000000, 1437955200000, 1438041600000, 1438128000000, 1438214400000, 1438300800000, 1438560000000, 1438646400000, 1438732800000, 1438819200000, 1438905600000, 1439164800000, 1439251200000, 1439337600000, 1439424000000, 1439510400000, 1439769600000, 1439856000000, 1439942400000, 1440028800000, 1440115200000, 1440374400000, 1440460800000, 1440547200000, 1440633600000, 1440720000000, 1440979200000, 1441065600000, 1441152000000, 1441238400000, 1441324800000, 1441670400000, 1441756800000, 1441843200000, 1441929600000, 1442188800000, 1442275200000, 1442361600000, 1442448000000, 1442534400000, 1442793600000, 1442880000000, 1442966400000, 1443052800000, 1443139200000, 1443398400000, 1443484800000, 1443571200000, 1443657600000, 1443744000000, 1444003200000, 1444089600000, 1444176000000, 1444262400000, 1444348800000, 1444608000000, 1444694400000, 1444780800000, 1444867200000, 1444953600000, 1445212800000, 1445299200000, 1445385600000, 1445472000000, 1445558400000, 1445817600000, 1445904000000, 1445990400000, 1446076800000, 1446163200000, 1446422400000, 1446508800000, 1446595200000, 1446681600000, 1446768000000, 1447027200000, 1447113600000, 1447200000000, 1447286400000, 1447372800000, 1447632000000, 1447718400000, 1447804800000, 1447891200000, 1447977600000, 1448236800000, 1448323200000, 1448409600000, 1448582400000, 1448841600000, 1448928000000, 1449014400000, 1449100800000, 1449187200000, 1449446400000, 1449532800000, 1449619200000, 1449705600000, 1449792000000, 1450051200000, 1450137600000, 1450224000000, 1450310400000, 1450396800000, 1450656000000, 1450742400000, 1450828800000, 1450915200000, 1451260800000, 1451347200000, 1451433600000, 1451520000000, 1451865600000, 1451952000000, 1452038400000, 1452124800000, 1452211200000, 1452470400000, 1452556800000, 1452643200000, 1452729600000, 1452816000000, 1453161600000, 1453248000000, 1453334400000, 1453420800000, 1453680000000, 1453766400000, 1453852800000, 1453939200000, 1454025600000, 1454284800000, 1454371200000, 1454457600000, 1454544000000, 1454630400000, 1454889600000, 1454976000000, 1455062400000, 1455148800000, 1455235200000, 1455580800000, 1455667200000, 1455753600000, 1455840000000, 1456099200000, 1456185600000, 1456272000000, 1456358400000, 1456444800000, 1456704000000, 1456790400000, 1456876800000, 1456963200000, 1457049600000, 1457308800000, 1457395200000, 1457481600000, 1457568000000, 1457654400000, 1457913600000, 1458000000000, 1458086400000, 1458172800000, 1458259200000, 1458518400000, 1458604800000, 1458691200000, 1458777600000, 1459123200000, 1459209600000, 1459296000000, 1459382400000, 1459468800000, 1459728000000, 1459814400000, 1459900800000, 1459987200000, 1460073600000, 1460332800000, 1460419200000, 1460505600000, 1460592000000, 1460678400000, 1460937600000, 1461024000000, 1461110400000, 1461196800000, 1461283200000, 1461542400000, 1461628800000, 1461715200000, 1461801600000, 1461888000000, 1462147200000, 1462233600000, 1462320000000, 1462406400000, 1462492800000, 1462752000000, 1462838400000, 1462924800000, 1463011200000, 1463097600000, 1463356800000, 1463443200000, 1463529600000, 1463616000000, 1463702400000, 1463961600000, 1464048000000, 1464134400000, 1464220800000, 1464307200000, 1464652800000, 1464739200000, 1464825600000, 1464912000000, 1465171200000, 1465257600000, 1465344000000, 1465430400000, 1465516800000, 1465776000000, 1465862400000, 1465948800000, 1466035200000, 1466121600000, 1466380800000, 1466467200000, 1466553600000, 1466640000000, 1466726400000, 1466985600000, 1467072000000, 1467158400000, 1467244800000, 1467331200000, 1467676800000, 1467763200000, 1467849600000, 1467936000000, 1468195200000, 1468281600000, 1468368000000, 1468454400000, 1468540800000, 1468800000000, 1468886400000, 1468972800000, 1469059200000, 1469145600000, 1469404800000, 1469491200000, 1469577600000, 1469664000000, 1469750400000, 1470009600000, 1470096000000, 1470182400000, 1470268800000, 1470355200000, 1470614400000, 1470700800000, 1470787200000, 1470873600000, 1470960000000, 1471219200000, 1471305600000, 1471392000000, 1471478400000, 1471564800000, 1471824000000, 1471910400000, 1471996800000, 1472083200000, 1472169600000, 1472428800000, 1472515200000, 1472601600000, 1472688000000, 1472774400000, 1473120000000, 1473206400000, 1473292800000, 1473379200000, 1473638400000, 1473724800000, 1473811200000, 1473897600000, 1473984000000, 1474243200000, 1474329600000, 1474416000000, 1474502400000, 1474588800000, 1474848000000, 1474934400000, 1475020800000, 1475107200000, 1475193600000, 1475452800000, 1475539200000, 1475625600000, 1475712000000, 1475798400000, 1476057600000, 1476144000000, 1476230400000, 1476316800000, 1476403200000, 1476662400000, 1476748800000, 1476835200000, 1476921600000, 1477008000000, 1477267200000, 1477353600000, 1477440000000, 1477526400000, 1477612800000, 1477872000000, 1477958400000, 1478044800000, 1478131200000, 1478217600000, 1478476800000, 1478563200000, 1478649600000, 1478736000000, 1478822400000, 1479081600000, 1479168000000, 1479254400000, 1479340800000, 1479427200000, 1479686400000, 1479772800000, 1479859200000, 1480032000000, 1480291200000, 1480377600000, 1480464000000, 1480550400000, 1480636800000, 1480896000000, 1480982400000, 1481068800000, 1481155200000, 1481241600000, 1481500800000, 1481587200000, 1481673600000, 1481760000000, 1481846400000, 1482105600000, 1482192000000, 1482278400000, 1482364800000, 1482451200000, 1482796800000, 1482883200000, 1482969600000, 1483056000000, 1483401600000, 1483488000000, 1483574400000, 1483660800000, 1483920000000, 1484006400000, 1484092800000, 1484179200000, 1484265600000, 1484611200000, 1484697600000, 1484784000000, 1484870400000, 1485129600000, 1485216000000, 1485302400000, 1485388800000, 1485475200000, 1485734400000, 1485820800000, 1485907200000, 1485993600000, 1486080000000, 1486339200000, 1486425600000, 1486512000000, 1486598400000, 1486684800000, 1486944000000, 1487030400000, 1487116800000, 1487203200000, 1487289600000, 1487635200000, 1487721600000, 1487808000000, 1487894400000, 1488153600000, 1488240000000, 1488326400000, 1488412800000, 1488499200000, 1488758400000, 1488844800000, 1488931200000, 1489017600000, 1489104000000, 1489363200000, 1489449600000, 1489536000000, 1489622400000, 1489708800000, 1489968000000, 1490054400000, 1490140800000, 1490227200000, 1490313600000, 1490572800000, 1490659200000, 1490745600000, 1490832000000, 1490918400000, 1491177600000, 1491264000000, 1491350400000, 1491436800000, 1491523200000, 1491782400000, 1491868800000, 1491955200000, 1492041600000, 1492387200000, 1492473600000, 1492560000000, 1492646400000, 1492732800000, 1492992000000, 1493078400000, 1493164800000, 1493251200000, 1493337600000, 1493596800000, 1493683200000, 1493769600000, 1493856000000, 1493942400000, 1494201600000, 1494288000000, 1494374400000, 1494460800000, 1494547200000, 1494806400000, 1494892800000, 1494979200000, 1495065600000, 1495152000000, 1495411200000, 1495497600000, 1495584000000, 1495670400000, 1495756800000, 1496102400000, 1496188800000, 1496275200000, 1496361600000, 1496620800000, 1496707200000, 1496793600000, 1496880000000, 1496966400000, 1497225600000, 1497312000000, 1497398400000, 1497484800000, 1497571200000, 1497830400000, 1497916800000, 1498003200000, 1498089600000, 1498176000000, 1498435200000, 1498521600000, 1498608000000, 1498694400000, 1498780800000, 1499040000000, 1499212800000, 1499299200000, 1499385600000, 1499644800000, 1499731200000, 1499817600000, 1499904000000, 1499990400000, 1500249600000, 1500336000000, 1500422400000, 1500508800000, 1500595200000, 1500854400000, 1500940800000, 1501027200000, 1501113600000, 1501200000000, 1501459200000, 1501545600000, 1501632000000, 1501718400000, 1501804800000, 1502064000000, 1502150400000, 1502236800000, 1502323200000, 1502409600000, 1502668800000, 1502755200000, 1502841600000, 1502928000000, 1503014400000, 1503273600000, 1503360000000, 1503446400000, 1503532800000, 1503619200000, 1503878400000, 1503964800000, 1504051200000, 1504137600000, 1504224000000, 1504569600000, 1504656000000, 1504742400000, 1504828800000, 1505088000000, 1505174400000, 1505260800000, 1505347200000, 1505433600000, 1505692800000, 1505779200000, 1505865600000, 1505952000000, 1506038400000, 1506297600000, 1506384000000, 1506470400000, 1506556800000, 1506643200000, 1506902400000, 1506988800000, 1507075200000, 1507161600000, 1507248000000, 1507507200000, 1507593600000, 1507680000000, 1507766400000, 1507852800000, 1508112000000, 1508198400000, 1508284800000, 1508371200000, 1508457600000, 1508716800000, 1508803200000, 1508889600000, 1508976000000, 1509062400000, 1509321600000, 1509408000000, 1509494400000, 1509580800000, 1509667200000, 1509926400000, 1510012800000, 1510099200000, 1510185600000, 1510272000000, 1510531200000, 1510617600000, 1510704000000, 1510790400000, 1510876800000, 1511136000000, 1511222400000, 1511308800000, 1511481600000, 1511740800000, 1511827200000, 1511913600000, 1512000000000, 1512086400000, 1512345600000, 1512432000000, 1512518400000, 1512604800000, 1512691200000, 1512950400000, 1513036800000, 1513123200000, 1513209600000, 1513296000000, 1513555200000, 1513641600000, 1513728000000, 1513814400000, 1513900800000, 1514246400000, 1514332800000, 1514419200000, 1514505600000, 1514851200000, 1514937600000, 1515024000000, 1515110400000, 1515369600000, 1515456000000, 1515542400000, 1515628800000, 1515715200000, 1516060800000, 1516147200000, 1516233600000, 1516320000000, 1516579200000, 1516665600000, 1516752000000, 1516838400000, 1516924800000, 1517184000000, 1517270400000, 1517356800000, 1517443200000, 1517529600000, 1517788800000, 1517875200000, 1517961600000}
var valList = []float64{27.260000, 27.405000, 27.370000, 27.370000, 27.610000, 27.400000, 27.290000, 27.815000, 26.810000, 28.230000, 30.130000, 29.455000, 30.370000, 31.250000, 30.900000, 31.550000, 31.865000, 31.310000, 31.250000, 32.485000, 32.295000, 33.000000, 32.560000, 32.925000, 34.020000, 33.115000, 33.940000, 34.165000, 33.750000, 34.000000, 34.135000, 33.495000, 33.630000, 33.850000, 33.485000, 33.720000, 33.265000, 32.475000, 32.120000, 34.215000, 34.340000, 35.375000, 35.365000, 34.680000, 33.525000, 32.277500, 32.225000, 32.310000, 32.895000, 32.325000, 32.905000, 33.125000, 33.710000, 34.245000, 33.965000, 34.060000, 33.785000, 33.260000, 33.420000, 33.655000, 34.055000, 33.985000, 33.910000, 33.740000, 33.345000, 33.415000, 34.120000, 34.265000, 33.990000, 35.155000, 36.215000, 35.620000, 34.905000, 35.810000, 36.355000, 36.000000, 36.100000, 35.840000, 35.250000, 35.365000, 35.120000, 34.905000, 35.525000, 36.310000, 35.260000, 35.510000, 34.655000, 35.345000, 35.235000, 35.700000, 36.115000, 35.105000, 34.685000, 33.645000, 34.965000, 35.520000, 35.885000, 35.180000, 35.840000, 35.400000, 35.540000, 36.090000, 35.725000, 36.010000, 35.995000, 36.205000, 35.455000, 35.415000, 35.145000, 34.790000, 35.100000, 36.475000, 36.670000, 36.495000, 36.225000, 37.505000, 38.285000, 38.740000, 37.990000, 38.540000, 38.440000, 38.670000, 38.540000, 38.525000, 38.580000, 38.015000, 37.825000, 37.945000, 37.675000, 37.710000, 36.970000, 37.550000, 37.230000, 37.000000, 37.620000, 37.680000, 38.730000, 39.000000, 38.440000, 38.740000, 39.520000, 39.350000, 39.770000, 39.720000, 39.650000, 38.500000, 38.980000, 39.000000, 38.760000, 38.910000, 38.040000, 37.950000, 37.740000, 38.070000, 38.470000, 38.070000, 37.510000, 37.270000, 36.620000, 36.430000, 37.010000, 36.670000, 37.320000, 36.650000, 36.040000, 35.760000, 35.850000, 35.850000, 35.090000, 36.190000, 36.080000, 36.360000, 36.660000, 36.730000, 37.550000, 36.620000, 36.250000, 36.180000, 35.680000, 34.840000, 35.950000, 37.040000, 37.080000, 36.470000, 35.430000, 35.240000, 35.220000, 34.860000, 35.000000, 34.000000, 32.730000, 33.660000, 33.600000, 33.290000, 33.660000, 33.290000, 33.600000, 32.770000, 33.950000, 34.330000, 34.890000, 35.280000, 34.680000, 34.270000, 34.470000, 34.570000, 34.500000, 34.480000, 34.130000, 36.260000, 36.540000, 36.900000, 36.270000, 35.810000, 36.490000, 36.190100, 35.770000, 36.040000, 37.150000, 38.600000, 38.540000, 38.560000, 38.970000, 38.560000, 39.020000, 38.670000, 38.570000, 38.170000, 38.150000, 38.400000, 38.730000, 38.930000, 37.430000, 37.570000, 36.550000, 37.390000, 37.870000, 38.280100, 37.950000, 39.590000, 40.060000, 39.580000, 39.550000, 39.540000, 39.620000, 40.560000, 39.370000, 40.550000, 40.120000, 41.729900, 40.510000, 40.300000, 39.130000, 38.820000, 39.470000, 38.360000, 38.910000, 39.440000, 39.550000, 40.100000, 38.140000, 36.050000, 35.590000, 35.430000, 35.460000, 35.430000, 34.880000, 35.030000, 34.720000, 35.000000, 35.330000, 35.000000, 35.230000, 35.000000, 35.090000, 34.620000, 34.130000, 33.280000, 33.620000, 33.380000, 33.320000, 32.880000, 32.850000, 32.640000, 32.310000, 33.640000, 33.800000, 34.080000, 34.300000, 35.100000, 35.710000, 34.500000, 34.130000, 34.510000, 33.480000, 32.350000, 32.630000, 32.730000, 33.180000, 34.030000, 34.780000, 35.030000, 35.960000, 36.950000, 37.650000, 38.450000, 39.000000, 38.860000, 39.260000, 38.999900, 38.700000, 38.760000, 38.870000, 37.920000, 37.250000, 36.930000, 37.360000, 37.420000, 37.000000, 37.150000, 36.520000, 36.700000, 36.000000, 36.660000, 36.170000, 36.000000, 36.620000, 36.100000, 36.340000, 36.210000, 35.790000, 35.980000, 36.140000, 36.320000, 36.020000, 35.650000, 34.720000, 35.090000, 35.190000, 34.840000, 34.620000, 34.890000, 34.960000, 35.140000, 34.750000, 34.340000, 33.890000, 34.340000, 33.990000, 34.200000, 34.150000, 34.220000, 35.230000, 34.720000, 34.020000, 34.630000, 34.420000, 34.600000, 34.100000, 34.400000, 34.260000, 33.500000, 33.500000, 33.300000, 33.400000, 33.220000, 34.450000, 32.870000, 32.390000, 32.890000, 33.560000, 33.000000, 32.780000, 32.780000, 34.170000, 33.600000, 33.700000, 33.560000, 34.380000, 34.120000, 33.560000, 33.190000, 33.210000, 33.150000, 33.340000, 32.830000, 33.370000, 33.170000, 33.120000, 33.670000, 33.340000, 32.880000, 33.310000, 33.910000, 33.490000, 33.840000, 33.490000, 33.870000, 33.610000, 33.550000, 33.460000, 33.280000, 32.960000, 33.160000, 33.900000, 33.830000, 33.310000, 32.960000, 32.490000, 30.970000, 31.950000, 31.370000, 31.970000, 32.900000, 32.640000, 31.540000, 31.520000, 31.160000, 31.120000, 30.960000, 31.010000, 29.820000, 29.250000, 28.750000, 27.950000, 29.670000, 31.710000, 30.830000, 31.380000, 31.640000, 30.670000, 32.020000, 30.480000, 30.730000, 31.910000, 31.410000, 30.660000, 31.440000, 30.410000, 31.170000, 31.510000, 32.820000, 33.800000, 32.860000, 33.280000, 32.970000, 32.920000, 33.580000, 33.430000, 33.670000, 33.990000, 35.230000, 34.230000, 34.600000, 33.970000, 32.130000, 32.780000, 32.900000, 31.890000, 31.070000, 31.160000, 30.260000, 29.660000, 29.850000, 29.410000, 29.050000, 30.230000, 29.360000, 30.510000, 32.830000, 32.620000, 30.880000, 30.740000, 30.450000, 30.470000, 30.110000, 29.690000, 29.430000, 29.470000, 29.540000, 28.320000, 28.490000, 28.280000, 28.840000, 29.530000, 29.600000, 29.560000, 30.040000, 28.850000, 28.750000, 29.140000, 29.090000, 28.760000, 28.490000, 28.560000, 27.900000, 27.210000, 26.290000, 26.860000, 27.700000, 27.100000, 26.850000, 26.870000, 26.970000, 27.260000, 27.740000, 28.560000, 28.240000, 27.920000, 27.780000, 26.820000, 27.630000, 27.640000, 27.750000, 28.410000, 29.750000, 29.020000, 29.000000, 29.160000, 29.440000, 28.520000, 28.230000, 27.840000, 27.350000, 26.950000, 28.090000, 27.600000, 27.920000, 28.880000, 27.660000, 28.050000, 28.660000, 28.650000, 28.530000, 29.120000, 29.720000, 28.820000, 28.900000, 28.970000, 29.480000, 29.830000, 30.870000, 30.610000, 30.960000, 30.830000, 31.600000, 31.270000, 30.960000, 31.660000, 32.590000, 32.740000, 32.770000, 32.900000, 32.810000, 32.870000, 33.380000, 32.940000, 32.800000, 33.290000, 33.820000, 33.790000, 34.120000, 34.520000, 34.280000, 34.030000, 34.600000, 35.060000, 35.010000, 35.220000, 35.150000, 35.180000, 35.380000, 35.000000, 34.980000, 35.110000, 34.820000, 34.730000, 34.590000, 34.570000, 33.920000, 33.910000, 34.210000, 34.080000, 34.000000, 33.850000, 33.840000, 34.490000, 34.260000, 34.050000, 33.760000, 33.650000, 33.880000, 34.350000, 34.220000, 33.390000, 33.240000, 33.130000, 32.600000, 32.340000, 31.800000, 31.670000, 31.640000, 31.510000, 31.240000, 30.390000, 30.370000, 30.340000, 30.080000, 29.800000, 29.490000, 30.380000, 29.700000, 29.550000, 29.430000, 28.750000, 28.900000, 28.270000, 27.870000, 25.510000, 25.890000, 26.470000, 26.570000, 26.370000, 25.870000, 25.750000, 25.940000, 25.830000, 25.250000, 25.230000, 25.090000, 25.160000, 26.080000, 26.720000, 27.000000, 26.810000, 26.570000, 26.450000, 25.750000, 24.810000, 23.090000, 22.950000, 21.760000, 22.170000, 23.010000, 23.000000, 22.960000, 23.340000, 23.390000, 23.100000, 23.100000, 23.000000, 22.560000, 22.580000, 22.990000, 22.340000, 22.890000, 23.670000, 23.300000, 23.730000, 23.500000, 24.000000, 23.410000, 24.330000, 23.550000, 22.780000, 21.750000, 22.140000, 21.000000, 22.180000, 22.930000, 23.960000, 23.250000, 23.240000, 23.320000, 22.730000, 22.340000, 22.970000, 23.840000, 23.000000, 23.140000, 23.300000, 22.430000, 21.050000, 20.910000, 20.470000, 21.250000, 21.840000, 21.580000, 21.530000, 22.360000, 22.180000, 22.260000, 22.340000, 22.250000, 22.440000, 22.590000, 20.490000, 20.290000, 20.770000, 22.370000, 21.120000, 20.260000, 20.410000, 19.640000, 19.510000, 19.020000, 19.110000, 19.190000, 18.780000, 18.320000, 18.070000, 16.840000, 16.540000, 16.240000, 17.590000, 17.220000, 17.240000, 15.910000, 15.860000, 15.160000, 15.430000, 15.290000, 15.700000, 16.160000, 17.010000, 17.680000, 17.400000, 17.810000, 17.200000, 17.090000, 17.670000, 17.660000, 17.020000, 15.710000, 16.780000, 17.080000, 16.320000, 15.710000, 15.850000, 15.470000, 16.530000, 16.050000, 16.950000, 19.360000, 18.940000, 18.810000, 19.270000, 20.400000, 19.680000, 20.240000, 20.170000, 20.820000, 21.070000, 20.790000, 20.310000, 20.420000, 19.490000, 19.270000, 20.210000, 20.090000, 19.410000, 20.530000, 19.350000, 20.440000, 20.330000, 18.920000, 19.710000, 20.130000, 20.350000, 20.390000, 20.860000, 21.890000, 22.520000, 21.440000, 22.020000, 21.860000, 20.680000, 22.410000, 21.570000, 21.870000, 22.600000, 23.250000, 22.890000, 22.640000, 21.800000, 21.630000, 21.020000, 22.000000, 21.670000, 22.880000, 22.740000, 22.260000, 22.610000, 22.420000, 22.350000, 22.300000, 23.710000, 23.490000, 23.130000, 23.780000, 23.140000, 22.600000, 22.390000, 23.700000, 23.400000, 23.460000, 23.000000, 22.880000, 22.410000, 23.080000, 23.790000, 23.710000, 23.330000, 24.030000, 24.290000, 24.370000, 23.790000, 23.650000, 24.140000, 24.370000, 24.800000, 24.190000, 24.110000, 23.910000, 23.780000, 22.810000, 23.380000, 22.890000, 23.430000, 23.570000, 24.230000, 23.690000, 24.020000, 23.820000, 24.360000, 24.790000, 24.650000, 24.610000, 25.700000, 24.850000, 25.560000, 24.920000, 25.020000, 24.630000, 24.260000, 24.660000, 25.040000, 24.980000, 25.520000, 25.390000, 24.530000, 24.530000, 24.560000, 26.280000, 25.330000, 25.810000, 25.260000, 25.280000, 26.160000, 25.850000, 25.500000, 26.000000, 25.960000, 26.140000, 25.450000, 25.080000, 25.370000, 24.690000, 24.500000, 24.410000, 24.200000, 23.410000, 23.770000, 23.600000, 23.560000, 24.580000, 23.770000, 24.050000, 24.370000, 24.460000, 24.740000, 25.050000, 24.330000, 24.280000, 24.720000, 25.200000, 25.330000, 25.100000, 25.120000, 25.950000, 25.710000, 25.830000, 25.940000, 25.620000, 25.510000, 25.310000, 25.320000, 24.700000, 24.440000, 24.690000, 24.390000, 24.800000, 25.500000, 26.220000, 25.920000, 26.020000, 25.110000, 25.150000, 24.710000, 24.520000, 24.500000, 24.790000, 25.550000, 25.720000, 25.540000, 24.820000, 24.590000, 25.420000, 25.580000, 25.890000, 25.600000, 25.540000, 25.750000, 25.550000, 25.570000, 24.030000, 23.550000, 23.710000, 23.040000, 22.330000, 22.160000, 22.220000, 21.090000, 21.890000, 21.250000, 20.940000, 20.370000, 20.920000, 20.700000, 20.460000, 20.990000, 20.270000, 20.300000, 20.140000, 20.450000, 20.330000, 20.310000, 21.690000, 21.320000, 21.310000, 22.160000, 22.250000, 21.880000, 21.780000, 22.360000, 23.140000, 22.370000, 23.040000, 23.180000, 22.300000, 24.440000, 23.060000, 22.720000, 23.200000, 23.870000, 23.740000, 23.830000, 24.640000, 24.170000, 23.480000, 23.310000, 22.520000, 22.650000, 21.850000, 21.940000, 22.190000, 22.360000, 22.400000, 22.890000, 23.240000, 22.800000, 22.810000, 22.980000, 21.930000, 22.400000, 22.740000, 22.800000, 22.410000, 22.280000, 22.300000, 22.210000, 21.850000, 21.530000, 21.730000, 22.720000, 22.320000, 22.350000, 23.000000, 23.080000, 23.080000, 22.320000, 21.180000, 21.640000, 20.660000, 21.280000, 23.750000, 23.730000, 23.340000, 23.870000, 24.170000, 23.900000, 23.860000, 24.140000, 23.920000, 23.010000, 22.930000, 22.840000, 22.950000, 22.600000, 22.440000, 22.040000, 22.110000, 22.520000, 22.460000, 22.680000, 23.120000, 22.820000, 22.590000, 23.040000, 22.780000, 22.630000, 22.540000, 22.670000, 22.390000, 22.300000, 22.750000, 22.470000, 22.520000, 22.880000, 22.840000, 23.190000, 23.630000, 24.190000, 23.470000, 23.950000, 24.050000, 24.880000, 24.500000, 24.920000, 24.640000, 24.510000, 24.540000, 24.470000, 24.000000, 23.940000, 24.240000, 24.070000, 24.390000, 24.210000, 23.790000, 24.170000, 24.110000, 24.660000, 23.350000, 23.850000, 24.020000, 24.070000, 23.650000, 23.940000, 23.750000, 23.990000, 25.250000, 24.890000, 24.610000, 24.260000, 23.720000, 23.100000, 23.130000, 23.610000, 23.500000, 23.500000, 22.730000, 22.840000, 22.740000, 22.240000, 22.280000, 22.070000, 21.400000, 21.740000, 21.860000, 21.950000, 22.360000, 23.110000, 22.880000, 23.130000, 23.000000, 23.460000, 23.040000, 22.280000, 22.300000, 22.430000, 22.880000, 23.860000, 24.480000, 24.580000, 24.620000, 25.160000, 25.090000, 24.910000, 24.950000, 24.160000, 24.260000, 24.640000, 25.060000, 25.060000, 25.430000, 25.370000, 25.530000, 24.880000, 25.330000, 24.980000, 24.900000, 25.120000, 25.230000, 25.180000, 25.100000, 25.320000, 24.570000, 25.190000, 24.670000, 24.570000, 24.640000, 24.280000, 24.750000, 24.890000, 24.480000, 24.720000, 24.380000, 24.560000, 24.490000, 24.510000, 24.340000, 24.590000, 24.860000, 24.890000, 25.030000, 24.770000, 24.550000, 24.720000, 24.940000, 25.550000, 26.660000, 27.060000, 26.900000, 26.520000, 26.120000, 26.150000, 26.570000, 27.130000, 26.750000, 26.780000, 26.460000, 26.340000, 26.130000, 26.340000, 26.730000, 26.810000, 27.000000, 26.830000, 26.590000, 26.480000, 26.720000, 26.740000, 26.570000, 26.190000, 25.710000, 25.610000, 25.440000, 25.420000, 26.230000, 25.990000, 25.690000, 25.780000, 25.500000, 25.520000, 25.670000, 25.270000, 24.790000, 24.410000, 25.200000, 26.850000, 27.400000, 27.750000, 27.970000, 28.230000, 28.600000, 27.450000, 27.800000, 28.250000, 28.310000, 28.890000, 28.810000, 28.360000, 28.580000, 28.910000, 29.240000, 29.130000, 29.280000, 29.490000, 29.400000, 29.250000, 29.390000, 29.210000, 29.230000, 28.510000, 28.430000, 28.280000, 27.950000, 27.890000, 28.030000, 28.390000, 27.580000, 27.680000, 27.450000, 26.870000, 27.130000, 26.760000, 26.830000, 27.700000, 27.750000, 28.010000, 28.320000, 28.690000, 28.310000, 29.180000, 29.090000, 28.580000, 28.860000, 28.940000, 29.140000, 28.190000, 28.300000, 28.780000, 27.800000, 27.880000, 27.570000, 27.710000, 28.170000, 28.600000, 28.270000, 27.860000, 27.560000, 26.870000, 26.750000, 26.290000, 25.820000, 24.760000, 23.650000, 24.540000}

func RealSample(num int) []sample {
	samples := []sample{}
	for i := 0; i < len(timeList) && i < num; i++ {
		samples = append(samples, sample{t: timeList[i], v: valList[i]})
	}
	return samples
}
