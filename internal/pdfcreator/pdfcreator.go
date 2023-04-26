package pdfcreator

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/viggneshvn/reddotstudios_contracts_backend/internal/contract"
)

type Pdf struct {
	pdf.Maroto
}

func CreateContractsPage(details *contract.Contract) (*string, error) {
	contractsPage := pdf.NewMarotoCustomSize(consts.Portrait, "Letter", "mm", 215.9, 210.0)

	contractsPage.SetBorder(true)

	contractsPage.Row(6, func() {
		contractsPage.Col(12, func() {
			// Note: Left-alignment is the default for text
			contractsPage.Text("Event Details", props.Text{
				Top:    1,
				Family: consts.Arial,
				Align:  consts.Center,
			})
		})
	})

	contractsPage.Row(20, func() {
		contractsPage.Col(12, func() {
			// Note: Left-alignment is the default for text
			contractsPage.Text(fmt.Sprintf("Event : %s", details.EventDetails.EventName), props.Text{
				Top:    2,
				Left:   2,
				Right:  2,
				Family: consts.Arial,
				Align:  consts.Left,
			})
			contractsPage.Text(fmt.Sprintf("Event Date: %s", details.EventDetails.EventDate), props.Text{
				Top:    6,
				Left:   2,
				Right:  2,
				Family: consts.Arial,
				Align:  consts.Left,
			})
			contractsPage.Text(fmt.Sprintf("Event Coverage Time: %s", details.EventDetails.EventCoverageTime), props.Text{
				Top:    10,
				Left:   2,
				Right:  2,
				Family: consts.Arial,
				Align:  consts.Left,
			})
			contractsPage.Text(fmt.Sprintf("Event Venue: %s", details.EventDetails.EventVenue), props.Text{
				Top:    14,
				Left:   2,
				Right:  2,
				Family: consts.Arial,
				Align:  consts.Left,
			})
		})
	})

	contractsPage.Row(10, func() { contractsPage.Text("") })

	contractsPage.Row(6, func() {
		contractsPage.Col(12, func() {
			contractsPage.Text("Deliverables", props.Text{
				Top:    1,
				Family: consts.Arial,
				Align:  consts.Center,
			})
		})
	})

	contractsPage.Row(6, func() {
		contractsPage.Col(1, func() {
			contractsPage.Text("S.No", props.Text{Top: 1, Align: consts.Center})
		})
		contractsPage.Col(3, func() {
			contractsPage.Text("Deliverable", props.Text{Top: 1, Align: consts.Center})
		})
		contractsPage.Col(2, func() {
			contractsPage.Text("Quantity", props.Text{Top: 1, Align: consts.Center})
		})
		contractsPage.Col(3, func() {
			contractsPage.Text("Mode", props.Text{Top: 1, Align: consts.Center})
		})
		contractsPage.Col(3, func() {
			contractsPage.Text("Delivered on or before**", props.Text{Top: 1, Align: consts.Center})
		})
	})

	for i, deliverable := range details.DeliverableDetails {
		contractsPage.Row(15, func() {
			contractsPage.Col(1, func() {
				contractsPage.Text(fmt.Sprintf("%d.", i+1), props.Text{Top: 1, Align: consts.Center})
			})
			contractsPage.Col(3, func() {
				desc := deliverable.Description
				if i == 0 {
					desc = fmt.Sprintf("%s*", desc)
				}
				contractsPage.Text(desc, props.Text{Top: 1, Align: consts.Center})
			})
			contractsPage.Col(2, func() {
				contractsPage.Text(deliverable.Quantity, props.Text{Top: 1, Align: consts.Center})
			})
			contractsPage.Col(3, func() {
				contractsPage.Text(deliverable.Mode, props.Text{Top: 1, Align: consts.Center})
			})
			contractsPage.Col(3, func() {
				contractsPage.Text(deliverable.DeliveryDate, props.Text{Top: 1, Align: consts.Center})
			})
		})
	}

	contractsPage.SetBorder(false)

	Astrisks := `Red Dot Studios does not provide RAW images/ video files unless specifically mentioned above in the section 2. Acquiring RAW images/ video comes at an additional cost.`

	contractsPage.Row(7, func() {
		contractsPage.Col(12, func() {
			contractsPage.Text("*", props.Text{
				Top:    1,
				Family: consts.Arial,
				Size:   8,
			})
			contractsPage.Text(Astrisks, props.Text{
				Family: consts.Arial,
				Top:    1,
				Left:   4,
				Size:   8,
			})
		})
	})

	Astrisks = `Red Dot Studios timelines for the delivery of projects depends on various factors which include the scale of the event, type of the service and number of deliverables, editing work according to the clients needs etc. We strive to deliver the first digital copy for the events by the date mentioned above. That being said, the delivery for event projects could take up to 2 months and wedding projects could take up to 6 months in special cases.`

	contractsPage.Row(8, func() {
		contractsPage.Col(12, func() {
			contractsPage.Text("**", props.Text{
				Family: consts.Arial,
				Size:   8,
			})
			contractsPage.Text(Astrisks, props.Text{
				Family: consts.Arial,
				Left:   4,
				Size:   8,
			})
		})
	})

	contractsPage.Row(10, func() { contractsPage.Text("") })

	contractsPage.SetBorder(true)

	contractsPage.Row(6, func() {
		contractsPage.Col(12, func() {
			contractsPage.Text("Payment Details", props.Text{
				Top:    1,
				Family: consts.Arial,
				Align:  consts.Center,
			})
		})
	})

	contractsPage.Row(6, func() {
		contractsPage.Col(1, func() {
			contractsPage.Text("S.No", props.Text{Top: 1, Align: consts.Center})
		})
		contractsPage.Col(4, func() {
			contractsPage.Text("Description", props.Text{Top: 1, Align: consts.Center})
		})
		contractsPage.Col(2, func() {
			contractsPage.Text("Amount", props.Text{Top: 1, Align: consts.Center})
		})
		contractsPage.Col(3, func() {
			contractsPage.Text("Mode", props.Text{Top: 1, Align: consts.Center})
		})
		contractsPage.Col(2, func() {
			contractsPage.Text("Status", props.Text{Top: 1, Align: consts.Center})
		})
	})

	sNo := 0
	if details.PaymentDetails.AdvancePaid != 0 {
		sNo++
		contractsPage.Row(6, func() {
			contractsPage.Col(1, func() {
				contractsPage.Text(fmt.Sprintf("%d.", sNo), props.Text{Top: 1, Align: consts.Center})
			})
			contractsPage.Col(4, func() {
				contractsPage.Text("Advance Paid", props.Text{Top: 1, Align: consts.Center})
			})
			contractsPage.Col(2, func() {
				contractsPage.Text(fmt.Sprintf("$%d", details.PaymentDetails.AdvancePaid), props.Text{Top: 1, Align: consts.Center})
			})
			contractsPage.Col(3, func() {
				contractsPage.Text(details.PaymentDetails.AdvancePaymentMode, props.Text{Top: 1, Align: consts.Center})
			})
			contractsPage.Col(2, func() {
				contractsPage.Text("Paid", props.Text{Top: 1, Align: consts.Center})
			})
		})
	}

	bookingFeeToBePaid := int64(math.Round(float64(details.PaymentDetails.TotalAmount) * 0.25))

	if bookingFeeToBePaid%10 != 0 {
		bookingFeeToBePaid = (((bookingFeeToBePaid / 10) + 1) * 10) - details.PaymentDetails.AdvancePaid
	} else {
		bookingFeeToBePaid = bookingFeeToBePaid - details.PaymentDetails.AdvancePaid
	}

	remainingProjectPayment := details.PaymentDetails.TotalAmount - bookingFeeToBePaid - details.PaymentDetails.AdvancePaid
	sNo++
	contractsPage.Row(10, func() {
		contractsPage.Col(1, func() {
			contractsPage.Text(fmt.Sprintf("%d.", sNo), props.Text{Top: 1, Align: consts.Center})
		})
		contractsPage.Col(4, func() {
			contractsPage.Text("Booking fee* - 25% (Non-refundable)", props.Text{Top: 2, Align: consts.Center})
		})
		contractsPage.Col(2, func() {
			contractsPage.Text(fmt.Sprintf("$%d", bookingFeeToBePaid), props.Text{Top: 2, Align: consts.Center})
		})
		contractsPage.Col(3, func() {
			contractsPage.Text("Venmo - @svarikut", props.Text{Top: 0.5, Align: consts.Center})
			contractsPage.Text("Zelle - 847-922-7653", props.Text{Top: 4, Align: consts.Center})
		})
		contractsPage.Col(2, func() {
			contractsPage.Text("To be paid", props.Text{Top: 2, Align: consts.Center})
		})
	})

	sNo++
	contractsPage.Row(6, func() {
		contractsPage.Col(1, func() {
			contractsPage.Text(fmt.Sprintf("%d.", sNo), props.Text{Top: 1, Align: consts.Center})
		})
		contractsPage.Col(4, func() {
			contractsPage.Text("Remaining Project Payment** - 75%", props.Text{Top: 1, Align: consts.Center})
		})
		contractsPage.Col(2, func() {
			contractsPage.Text(fmt.Sprintf("$%d", remainingProjectPayment), props.Text{Top: 1, Align: consts.Center})
		})
		contractsPage.Col(3, func() {
			contractsPage.Text("Cash Only", props.Text{Top: 1, Align: consts.Center})
		})
		contractsPage.Col(2, func() {
			contractsPage.Text("To be paid", props.Text{Top: 1, Align: consts.Center})
		})
	})

	contractsPage.Row(6, func() {
		contractsPage.Col(1, func() {
			contractsPage.Text("", props.Text{Top: 1, Align: consts.Center})
		})
		contractsPage.Col(4, func() {
			contractsPage.Text("Total - 100%", props.Text{Top: 1, Align: consts.Center})
		})
		contractsPage.Col(2, func() {
			contractsPage.Text(fmt.Sprintf("$%d", details.PaymentDetails.TotalAmount), props.Text{Top: 1, Align: consts.Center})
		})
		contractsPage.Col(3, func() {
			contractsPage.Text("", props.Text{Top: 1, Align: consts.Center})
		})
		contractsPage.Col(2, func() {
			contractsPage.Text("", props.Text{Top: 1, Align: consts.Center})
		})
	})

	contractsPage.SetBorder(false)

	Astrisks = `Booking fee(25%) has to be paid during the time of signing this contract. We do not guarantee the availability of our team for the event date until this payment is made in full.`

	contractsPage.Row(7, func() {
		contractsPage.Col(12, func() {
			contractsPage.Text("*", props.Text{
				Top:    1,
				Family: consts.Arial,
				Size:   8,
			})
			contractsPage.Text(Astrisks, props.Text{
				Family: consts.Arial,
				Top:    1,
				Left:   4,
				Size:   8,
			})
		})
	})

	Astrisks = `Remaining Project Payment(75%) has to be paid on the day of the event in cash. We do not accept any other mode of payment except cash; there is no exception to this policy. Editing work only begins on the receipt of complete payment.`

	contractsPage.Row(8, func() {
		contractsPage.Col(12, func() {
			contractsPage.Text("**", props.Text{
				Family: consts.Arial,
				Size:   8,
			})
			contractsPage.Text(Astrisks, props.Text{
				Family: consts.Arial,
				Left:   4,
				Size:   8,
			})
		})
	})

	subDir := "contracts"
	_, err := os.Stat(subDir)
	if os.IsNotExist(err) {
		err := os.Mkdir(subDir, 0755)
		if err != nil {
			return nil, fmt.Errorf("could not create contracts directory with error : %w", err)
		}
	}

	fileName := fmt.Sprintf("contracts/contract-specifics-%s-%s", strings.ReplaceAll(details.EventDetails.EventName, " ", "_"), strings.ReplaceAll(details.ClientDetails.ClientName, " ", "_"))
	err = contractsPage.OutputFileAndClose(fileName + ".pdf")
	if err != nil {
		return nil, fmt.Errorf("could not save pdf file with error : %w", err)
	}

	return &fileName, nil
}

func CreateTermsPage(details *contract.Contract) (*string, error) {
	termsPage := pdf.NewMarotoCustomSize(consts.Portrait, "Letter", "mm", 215.9, 138.0)
	termsPage.SetPageMargins(2, 3, 5)

	termsPage.Row(28, func() {
		termsPage.Col(1, func() {
			termsPage.Text("1.", props.Text{Top: 1, Align: consts.Center})
		})
		termsPage.Col(11, func() {
			termsPage.Text("Reschedule Policy:", props.Text{Top: 1, Family: consts.Arial,
				Style: consts.Bold})
			termsPage.Text("We understand that event dates and times can change due to several factors. We ", props.Text{Top: 1, Left: 34, Family: consts.Arial})
			termsPage.Text(fmt.Sprintf("accommodate up to 1 hour of delay/ prepone in the event time on the day of the event if informed 4 hours prior to the start of the event time. The Red Dot Studio team will also try to accommodate requests to extend stay to cover the event if the event runs longer than anticipated. That being this request will come at an extra hourly prorated cost of $%d/hr which is non-negotiable and is subject to availability. Our schedules are packed during busy months and the team might have to cover an event before or after the client’s event. We always encourage our clients to book our time conservatively if they anticipate any delays. Red Dot Studios allows one reschedule of the event if informed 24 hrs prior to the date of the event provided project payment is made in full while requesting the reschedule. There are no exceptions to this clause.", details.PaymentDetails.PerHourExtra), props.Text{Top: 4.5, Family: consts.Arial})
		})
	})

	termsPage.Row(18, func() {
		termsPage.Col(1, func() {
			termsPage.Text("2.", props.Text{Top: 1, Align: consts.Center})
		})
		termsPage.Col(11, func() {
			termsPage.Text("Cancellation/ Termination:", props.Text{Top: 1, Family: consts.Arial,
				Style: consts.Bold})
			termsPage.Text("Client may decide to terminate this agreement at any time upon a written", props.Text{Top: 1, Left: 46, Family: consts.Arial})
			termsPage.Text("notification(Email, whatsapp, instagram) to Red Dot Studios. After a written notification, this agreement would be deemed void. Red Dot Studios shall be entitled to retain the booking advance made by the client. Red Dot Studios is entitled to take other bookings for the event date after the termination of the contract and any further requests will only be subject to availability and would require drafting a new contract.", props.Text{Top: 4.5, Family: consts.Arial})
		})
	})

	termsPage.Row(14, func() {
		termsPage.Col(1, func() {
			termsPage.Text("3.", props.Text{Top: 1, Align: consts.Center})
		})
		termsPage.Col(11, func() {
			termsPage.Text("Modifications to video deliverables:", props.Text{Top: 1, Family: consts.Arial,
				Style: consts.Bold})
			termsPage.Text("Client agrees to our creative choices and artistic/style decisions that", props.Text{Top: 1, Left: 61.5, Family: consts.Arial})
			termsPage.Text("we make during editing. Once we deliver the first digital copy we allow the client to request up-to two revisions both of which need to be requested within one week of the delivered digital copy. Final soft copy for the project will be delivered to the client after the second revision, and the project will be termed Completed.", props.Text{Top: 4.5, Family: consts.Arial})
		})
	})

	termsPage.Row(7, func() {
		termsPage.Col(1, func() {
			termsPage.Text("4.", props.Text{Top: 1, Align: consts.Center})
		})
		termsPage.Col(11, func() {
			termsPage.Text("Data Retention Policy:", props.Text{Top: 1, Family: consts.Arial,
				Style: consts.Bold})
			termsPage.Text("We erase all the client data after the completion of the project and do not take any", props.Text{Top: 1, Left: 39, Family: consts.Arial})
			termsPage.Text("additional requests for changes. ", props.Text{Top: 4.5, Family: consts.Arial})
		})
	})

	termsPage.Row(17.5, func() {
		termsPage.Col(1, func() {
			termsPage.Text("5.", props.Text{Top: 1, Align: consts.Center})
		})
		termsPage.Col(11, func() {
			termsPage.Text("Copyright:", props.Text{Top: 1, Family: consts.Arial,
				Style: consts.Bold})
			termsPage.Text("Red Dot Studios shall retain the copyright to all the photographs and/or videography shot during", props.Text{Top: 1, Left: 19, Family: consts.Arial})
			termsPage.Text("the event. The Client shall not remove or alter any watermarks, logos, or other identification marks included on the photographs and/or videography without the prior written consent. Red Dot Studios also holds the rights to use the edited videos and photos for the purpose of promoting our business in digital media, including but not limited to our website and social media.", props.Text{Top: 4.5, Family: consts.Arial})
		})
	})

	termsPage.Row(7.5, func() {
		termsPage.Col(1, func() {
			termsPage.Text("6.", props.Text{Top: 1, Align: consts.Center})
		})
		termsPage.Col(11, func() {
			termsPage.Text("Additional Services:", props.Text{Top: 1, Family: consts.Arial,
				Style: consts.Bold})
			termsPage.Text("The Client may request additional services from the Red Dot Studios, but such", props.Text{Top: 1, Left: 35, Family: consts.Arial})
			termsPage.Text("requests must be made before the event date and such requests will only be entertained subject to availability.", props.Text{Top: 4.5, Family: consts.Arial})
		})
	})

	termsPage.Row(7.5, func() {
		termsPage.Col(1, func() {
			termsPage.Text("7.", props.Text{Top: 1, Align: consts.Center})
		})
		termsPage.Col(11, func() {
			termsPage.Text("Limitation of Liability:", props.Text{Top: 1, Family: consts.Arial,
				Style: consts.Bold})
			termsPage.Text("The Client may request additional services from the Red Dot Studios, but such", props.Text{Top: 1, Left: 38, Family: consts.Arial})
			termsPage.Text("requests must be made before the event date and such requests will only be entertained subject to availability.", props.Text{Top: 4.5, Family: consts.Arial})
		})
	})

	termsPage.Row(7.5, func() {
		termsPage.Col(1, func() {
			termsPage.Text("8.", props.Text{Top: 1, Align: consts.Center})
		})
		termsPage.Col(11, func() {
			termsPage.Text("Entire Agreement:", props.Text{Top: 1, Family: consts.Arial,
				Style: consts.Bold})
			termsPage.Text("This Agreement constitutes the entire agreement between the parties and supersedes", props.Text{Top: 1, Left: 32, Family: consts.Arial})
			termsPage.Text("all prior negotiations, representations, understandings, and agreements between the parties.", props.Text{Top: 4.5, Family: consts.Arial})
		})
	})

	subDir := "contracts"
	_, err := os.Stat(subDir)
	if os.IsNotExist(err) {
		err := os.Mkdir(subDir, 0755)
		if err != nil {
			return nil, fmt.Errorf("could not create contracts directory with error : %w", err)
		}
	}

	fileNameT := fmt.Sprintf("contracts/contract-terms-%s-%s", strings.ReplaceAll(details.EventDetails.EventName, " ", "_"), strings.ReplaceAll(details.ClientDetails.ClientName, " ", "_"))
	// pdfName := fmt.Sprintf("contacts/contract-%s-%s.pdf", strings.ReplaceAll("House Warming", " ", "_"), strings.ReplaceAll("Sainath", " ", "_"))
	err = termsPage.OutputFileAndClose(fileNameT + ".pdf")
	if err != nil {
		return nil, fmt.Errorf("could not save pdf file with error : %w", err)
	}

	return &fileNameT, nil
}

func CleanUpPdfs() error {
	err := os.RemoveAll("contracts")
	if err != nil {
		return fmt.Errorf("failure while cleaning up pdfs from contracts directory with err : %w", err)
	}
	return nil
}
