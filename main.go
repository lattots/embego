package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/lattots/embego/pkg/util"
)

type Response struct {
	DocumentEmbedding  []float64   `json:"document_embedding"`
	SentenceEmbeddings [][]float64 `json:"sentence_embeddings"`
}

func main() {
	texts := []string{
		"Ortlieb Duffle 110 on vesitiivis, j\u00e4reist\u00e4 materiaaleista valmistettu 110L reissulaukku, jota voit kantaa my\u00f6s sel\u00e4ss\u00e4. Vesitiivis pitk\u00e4 vetoketju, joka avaa ison suuaukon helpottaen pakkaamista ja purkamista.\nIP 67 -standardin mukaisesti vesitiivis, my\u00f6s TIZIP-vetoketju\nJ\u00e4re\u00e4t ja kest\u00e4v\u00e4t PS620- ja PS620C-kankaat\nHitsatut saumat\nKanna kuten reppua, pehmustetut olkaimet\nOlkaimet voi s\u00e4\u00e4t\u00e4\u00e4 my\u00f6s lyhyiksi kantokahvoiksi\nSis\u00e4puoliset kompressiohihnat puristavat tavarat minimitilaan ja v\u00e4hent\u00e4v\u00e4t vetoketjulle kohdistuvaa rasitusta\nKaksi sis\u00e4puoilsta vetoketjutaskua\nUlkopuolinen verkkorakenteinen tasku (ei vedenpit\u00e4v\u00e4 :)\nKiinnikepisteet remmeille, tms. ylim\u00e4\u00e4r\u00e4isten varusteiden kiinnityst\u00e4 varten\nValmistettu Saksassa\nViiden vuoden takuu\nMitat\n34 \u00d770 \u00d7 46 cm\nTilavuus: 110 litraa\n\nPaino\n1490 g\n\nTuotteen koko pakattuna tai myyntipakkauksen koko on arviolta 500 x 360 x 120 mm.\n\n\nVedenpit\u00e4vyys\nKyll\u00e4\n\n\nPaino\n1,52 kg(Sis\u00e4lt\u00e4\u00e4 mahdollisen myyntipakkauksen painon)\n\n\nTilavuus\n110 litraa\n\n\nTakuu\n60 kk",
		"Hanwagin naisten miellytt\u00e4v\u00e4t vaelluskeng\u00e4t. Keve\u00e4 huoliteltu rakenne tuo askeliin mukavuutta ja Hanwagin EcoShell kalvo pit\u00e4\u00e4 jalkasi kuivina. Blueridge Low Lady ES on luottovalinta ulkoilukeng\u00e4ksi. \n\nModernin n\u00e4k\u00f6isiss\u00e4 jalkineissa on k\u00e4ytetty kest\u00e4v\u00e4\u00e4 Perwanger-nahkaa ja Global Recycling Standard -hyv\u00e4ksytty\u00e4 kierr\u00e4tetty\u00e4 polyamidia. Keng\u00e4t on tuotettu kokonaan Euroopassa. Hanwag on suunnitellut oman PU-kalvon, joka hydrofiilisella menetelm\u00e4ll\u00e4\u00e4n saavuttaa vedenpit\u00e4vyyden ja hyv\u00e4n hengitt\u00e4vyyden ilman fluorihiiilivety\u00e4. T\u00e4m\u00e4 kalvo on yhdistetty kolmikerroslaminaattiin. Mukavuutta on lis\u00e4\u00e4m\u00e4ss\u00e4 v\u00e4lipohja, joka on tehty kolmiulotteisella rakenteella. Ulkopohjassa on ekstra leve\u00e4ll\u00e4 kulutuspinnalla 4 mm syv\u00e4 kuviointi, joka parantaa pitoa. \nKest\u00e4v\u00e4 3-kerroksinen PFC-vapaa EcoShell PU -kalvo\nVakaa ja pehme\u00e4 v\u00e4lipohja, jossa on k\u00e4ytetty 3D rakennetta\nLaadukasta Perwanger-nahkaa ja kulutusta kest\u00e4v\u00e4\u00e4 polyamidia p\u00e4\u00e4llyskankaassa\nTerragrip Hike Pro -ulkopohja, jossa 4 mm syv\u00e4 kuviointi ekstra leve\u00e4ll\u00e4 kulutuspinnalla\n\nMateriaalit\nKalvo: EcoShell Footwear\nP\u00e4\u00e4llysmateriaali: Mokkanahka ja tekstiili\nUlkopohja: Hanwag Hike Pro\n\nPohjan j\u00e4ykkyys: AB\nNauhan pituus: UK 3,5 - 6 = 105cm | UK 6,5 - 9 = 115cm\n\nPaino: 820 g (yksi pari, koko UK 5)\n\nTuotteen koko pakattuna tai myyntipakkauksen koko on arviolta 300 x 150 x 200 mm.\n\n\nJ\u00e4ykkyys  \nA/B\n\n\nSukupuoli\nNaiset\n\n\nVedenpit\u00e4vyys\nGore-Tex vastaava\n\n\nVarren pituus\nMatala\n\n\nPaino\n0,82 kg\n(myyntipakkauksen kanssa 1 kg)\n\n\nTakuu\n24 kk\n\n",
		"Rab Sleep Limit: -35\u00b0C\n\nRab expedition 1200 Down Sleeping Bag on nyt sitten sit\u00e4 ehdottominta parasta talvipussilaatua.\n\nExpedition makuupussit on suunniteltu todella koviin olosuhteisiin ja jo aikamoiselle extremematkailijalle. Olit sitten matkalla kasitonniselle vuorelle tai tutkimusmatkalle arktiseen milj\u00f6\u00f6seen t\u00e4m\u00e4 pussi huolehtii l\u00e4mm\u00f6st\u00e4 niin hyvin kuin n\u00e4iss\u00e4 olosuhteissa se on mahdollista.\n\nExpedition 1200 tarjoaa l\u00e4mp\u00f6\u00e4 ja suojaa jopa -35 \u00b0C:n l\u00e4mp\u00f6tiloissa ja siin\u00e4 on 1200 g korkealaatuista untuvaa 850FP:t\u00e4. Erityisesti retkikuntak\u00e4ytt\u00f6\u00f6n suunniteltu Expedition 1200 on hieman ylimitoitettu, jotta isommat Expedition-vaatteetkin sopivat sis\u00e4\u00e4n. Pussissa k\u00e4ytet\u00e4\u00e4n korkealaatuisten Pertex\u00ae-sis\u00e4- ja p\u00e4\u00e4llyskankaiden yhdistelm\u00e4\u00e4.  Pussi hylkii hienosti vett\u00e4 ja huppu minimoi l\u00e4mm\u00f6n uloss\u00e4teilyn. Expeditionin makuupussit t\u00e4ytet\u00e4\u00e4n k\u00e4sin Derbyshiress\u00e4, Isossa-Britanniassa. T\u00e4ll\u00e4 tavalla untuvat s\u00e4ilyv\u00e4t korkealuokkaisina.\n\nPertex\u00ae Quantum Pro ulkokangas\nPertex\u00ae Quantum -sis\u00e4kangas Polygiene\u00ae Stays Fresh -teknologialla\nPertex\u00ae Quantum Pro sis\u00e4kangas hupussa ja jaloissa\n850FP eurooppalainen hanhenuntuva Nikwax-fluorihiilivapaalla hydrofobisella viimeistelyll\u00e4 (1200g  42.3oz)\nRab\u00ae Fluorocarbon free Hydrofobic Down kehitetty yhdess\u00e4 Nikwax\u00ae:in kanssa\n\"Muumion\" muotoinen\nYhteensopiva untuvavaatteilla (isompi sis\u00e4vuori kuin muissa Rab\u00ae-makuupusseissa)\n\u00be pituus YKK  vetoketju, est\u00e4\u00e4 jalkop\u00e4\u00e4n l\u00e4mp\u00f6vuotoja\nSis\u00e4ll\u00e4 l\u00e4mp\u00f6kaulus\nKompressiopussi\nK\u00e4sin t\u00e4ytetty Derbyshiress\u00e4 Englannissa\nMukana my\u00f6s puuvillainen s\u00e4ilytyspussi kotis\u00e4ilytyst\u00e4 varten\n\nMitat\nNukkujan pituus Regular: <195 cm\nNukkujan pituus Short: <170 cm\n\nRab nukkumisraja\n-35 \u00b0C\n\nPaino\n1840 g\n\nMateriaali\nUlkokangas: 100% nyloni \nSis\u00e4kangas: 100% nyloni\nEriste: 850FP hanhen untuva\n\nRAB kokotaulukko \n\n\nTuotteen koko pakattuna tai myyntipakkauksen koko on arviolta 370 x 280 x 280 mm.\n\n\nEriste T\u00e4yte\nUntuva\n\n\nPaino\n1,84 kg\n(myyntipakkauksen kanssa 2 kg)",
		"Rab Sleep Limit: -23\u00b0C (-10\u00b0F)\n\nRabin naisten talvimakuupussi kylmiin olosuhteisiin. T\u00e4m\u00e4 pussukka pit\u00e4\u00e4 varmasti l\u00e4mpim\u00e4n\u00e4. \n\nAndes Infinium 800 on t\u00e4ytetty 800 g:lla eritt\u00e4in korkeaa 800FP hanhenuntuvaa (parhaan l\u00e4mp\u00f6arvon omaavaa untuvaa) Fp 800 tai enemm\u00e4n-> laadultaan erinomainen. L\u00e4mp\u00f6\u00e4 vartaloa kohti heijastava Rab TILT -tekniikka n\u00e4kyy hopeisena kiiltona pussin sis\u00e4ll\u00e4. V\u00e4hent\u00e4\u00e4 l\u00e4mp\u00f6h\u00e4vikki\u00e4 jopa 32 %. TILT parantaa merkitt\u00e4v\u00e4sti l\u00e4mm\u00f6neristyst\u00e4 heikent\u00e4m\u00e4tt\u00e4 pakkauskokoa, hengitt\u00e4vyytt\u00e4 tai painoa.\n\nMuumioistuvuus tarjoaa saman l\u00e4mp\u00f6tason kuin Rabin Expedition-makuupussit, mutta pienemm\u00e4ss\u00e4 pakkauksessa. Naisellisen leikkauksen omaava makuupussi on hieman lyhyempi, olkap\u00e4iden tila on pienempi ja lantion kohdalla on enemm\u00e4n tilaa parantaen mukavuutta.\n\nGore-Tex Infinium\u2122 Windstopper\u00ae -ulkopuoli tarjoaa vankan, tuulenpit\u00e4v\u00e4n suojan j\u00e4isilt\u00e4 tuulilta. \u200b\u200b3D-kaulus kiristyy tiukasti kasvojen ymp\u00e4rille tiivist\u00e4m\u00e4\u00e4n l\u00e4mp\u00f6\u00e4. \u00be-pituisessa YKK\u00ae-vetoketjussa on oma sis\u00e4inen untuvilla t\u00e4ytetty \"v\u00e4lilevy\", joka est\u00e4\u00e4 l\u00e4mm\u00f6n karkaamisen.\n\nIsossa-Britanniassa k\u00e4sin t\u00e4ytetyt Rab-untuvapussit hy\u00f6tyv\u00e4t siit\u00e4, ett\u00e4 untuvia ei ole koskaan puristettu ennen k\u00e4ytt\u00f6\u00e4, mik\u00e4 s\u00e4ilytt\u00e4\u00e4 untuvan erinomaisen tason ja l\u00e4mm\u00f6n.\n800FP R.D.S. -sertifioitua eurooppalaista hanhenuntuvaa\nUntuva Rab Fluorihiilivapaata hydrofobista untuvaa, suunniteltu yhdess\u00e4 Nikwaxin kanssa\nGore-Tex Infinium Windstopper-ulkopuoli tarjoaa s\u00e4\u00e4suojan ja hyv\u00e4n hengitt\u00e4vyyden\nTILT-vuori auttaa heijastamaan ja v\u00e4hent\u00e4m\u00e4\u00e4n l\u00e4mp\u00f6h\u00e4vi\u00f6it\u00e4\nTrapezoidal boxwall -rakenne tuo entist\u00e4 enemm\u00e4n l\u00e4mp\u00f6\u00e4\n3D kaulus ei p\u00e4\u00e4st\u00e4 l\u00e4mp\u00f6\u00e4 karkaamaan\nGore-Tex Infinium -materiaali p\u00e4\u00e4n ja jalkojen kohdalla lis\u00e4ten veden hylkivyytt\u00e4\nYKK 3 / 4 -pituinen p\u00e4\u00e4vetoketju, jossa untuvilla t\u00e4ytetty v\u00e4lilevy l\u00e4mm\u00f6n ker\u00e4\u00e4miseksi\nYKK-vetoketjullinen pikkutasku sis\u00e4puolella\nMuotoiltu tila jaloille, jotta nukkuessa ne saavat olla rennosti\nMukana vetoketjullinen vett\u00e4hylkiv\u00e4 puuvillas\u00e4kki ja vedenpit\u00e4v\u00e4 pussukka\n\nTekniset tiedot \nPaino: 1360g / 47.9oz\nYmp\u00e4rysmitta olkap\u00e4iden kohdalta: 160cm / 63 inch\nNukkujan pituus: <180 cm\nKoko pakattuna: 45 X 28cm\n\nRAB kokotaulukko \n\n\nTuotteen koko pakattuna tai myyntipakkauksen koko on arviolta 470 x 300 x 300 mm.\n\n\nComfort\nEi ilmoitettu\n\n\nEriste/T\u00e4yte\nUntuva\n\n\nExtreme\nEi ilmoitettu\n\n\nLimit\n-23 \u00b0C\n\n\nPaino\n1,395 kg\n(myyntipakkauksen kanssa 1,49 kg)\n\n\nValmistusmaa\nIso-Britannia Iso-Britannia",
	}

	var embeddings []Response

	for _, text := range texts {
		body, err := json.Marshal(map[string]string{
			"text": text,
		})
		if err != nil {
			panic(err)
		}
		resp, err := http.Post("http://localhost:8081/embedding", "application/json", bytes.NewReader(body))
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		jsonString, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		var response Response
		err = json.Unmarshal(jsonString, &response)
		if err != nil {
			panic(err)
		}

		embeddings = append(embeddings, response)
	}

	for i := range embeddings {
		if i < len(embeddings)-1 {
			fmt.Println(util.CosineSimilarity(embeddings[i].DocumentEmbedding, embeddings[i+1].DocumentEmbedding))
		} else {
			fmt.Println(util.CosineSimilarity(embeddings[i].DocumentEmbedding, embeddings[0].DocumentEmbedding))
		}
	}
}
