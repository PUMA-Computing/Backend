package services

import (
	"Backend/configs"
	"Backend/internal/database/app"
	"context"
	"github.com/google/uuid"
	"github.com/mailgun/mailgun-go/v4"
	"log"
)

type MailgunService struct {
	domain      string
	apiKey      string
	senderEmail string
	mailgun     *mailgun.MailgunImpl
}

func NewMailgunService(domain, apiKey, senderEmail string) *MailgunService {
	return &MailgunService{
		domain:      domain,
		apiKey:      apiKey,
		senderEmail: senderEmail,
		mailgun:     mailgun.NewMailgun(domain, apiKey),
	}
}

type EmailData struct {
	Name             string
	VerificationLink string
	OTPCode          string
}

func (ms *MailgunService) SendOTPEmail(to, otpCode string) error {
	subject := "One Time Password"

	data := EmailData{
		OTPCode: otpCode,
	}

	// Construct the HTML email body
	body := `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office">
<head>
	<!--[if (gte mso 9)|(IE)]>
	<xml>
		<o:OfficeDocumentSettings>
			<o:AllowPNG/>
			<o:PixelsPerInch>96</o:PixelsPerInch>
		</o:OfficeDocumentSettings>
	</xml>
	<![endif]-->
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	<meta name="viewport" content="width=device-width" />
	<meta http-equiv="X-UA-Compatible" content="IE=edge" />
	<title>Notify - Notify Email Newsletter</title>

	<!-- Google Fonts Link -->
	<link href="https://fonts.googleapis.com/css?family=Open+Sans:300,400,600,700" rel="stylesheet" />
	<link href="https://fonts.googleapis.com/css?family=Lora:400,700" rel="stylesheet" />
	<style type="text/css">

		/*------ Client-Specific Style ------ */
		@-ms-viewport{width:device-width;}
		table, td{mso-table-lspace:0pt; mso-table-rspace:0pt;}
		img{-ms-interpolation-mode:bicubic; border: 0;}
		p, a, li, td, blockquote{mso-line-height-rule:exactly;}
		p, a, li, td, body, table, blockquote{-ms-text-size-adjust:100%; -webkit-text-size-adjust:100%;}
		#outlook a{padding:0;}
		.ReadMsgBody{width:100%;} .ExternalClass{width:100%;}
		.ExternalClass,.ExternalClass div,.ExternalClass font,.ExternalClass p,.ExternalClass span,.ExternalClass td,img{line-height:100%;}

		/*------ Reset Style ------ */
		*{-webkit-text-size-adjust:none;-webkit-text-resize:100%;text-resize:100%;}
		table{border-spacing: 0 !important;}
		h1, h2, h3, h4, h5, h6, p{display:block; Margin:0; padding:0;}
		img, a img{border:0; height:auto; outline:none; text-decoration:none;}
		#bodyTable, #bodyCell{ margin:0; padding:0; width:100%;}
		body {height:100%; margin:0; padding:0; width:100%;}

		.appleLinks a {color: #c2c2c2 !important; text-decoration: none;}
        span.preheader { display: none !important; }

		/*------ Google Font Style ------ */
		[style*="Open Sans"] {font-family:'Open Sans', Helvetica, Arial, sans-serif !important;}
		[style*="Lora"] {font-family:'Lora', Georgia, Times, serif !important;}

		/*------ General Style ------ */
		.wrapperWebview, .wrapperBody, .wrapperFooter{width:100%; max-width:600px; Margin:0 auto;}

		/*------ Column Layout Style ------ */
		.tableCard {text-align:center; font-size:0;}
		
		/*------ Images Style ------ */
		.imgHero img{ width:600px; height:auto; }
		
	</style>

	<style type="text/css">
		/*------ Media Width 480 ------ */
		@media screen and (max-width:640px) {
			table[class="wrapperWebview"]{width:100% !important; }
			table[class="wrapperEmailBody"]{width:100% !important; }
			table[class="wrapperFooter"]{width:100% !important; }
			td[class="imgHero"] img{ width:100% !important;}
			.hideOnMobile {display:none !important; width:0; overflow:hidden;}
		}
	</style>

</head>

<body style="background-color:#F9F9F9;">
<center>
		<table border="0" cellpadding="0" cellspacing="0" width="100%" style="table-layout:fixed;background-color:#F9F9F9;" id="bodyTable">
	<tbody><tr>
		<td align="center" valign="top" style="padding-right:10px;padding-left:10px;" id="bodyCell">
		<!--[if (gte mso 9)|(IE)]><table align="center" border="0" cellspacing="0" cellpadding="0" style="width:600px;" width="600"><tr><td align="center" valign="top"><![endif]-->

		<!-- Email Wrapper Webview Open //-->
		<table border="0" cellpadding="0" cellspacing="0" style="max-width:600px;" width="100%" class="wrapperWebview">
			<tbody><tr>
				<td align="center" valign="top">
					<!-- Content Table Open // -->
					<table border="0" cellpadding="0" cellspacing="0" width="100%">
						<tbody><tr>
							<td align="right" valign="middle" style="padding-top: 20px; padding-right: 0px;" class="webview">
								<!-- Email View in Browser // -->
								<a class="text hideOnMobile" href="https://computing.president.ac.id" target="_blank" style="color:#777777; font-family:'Open Sans', Helvetica, Arial, sans-serif; font-size:12px; font-weight:400; font-style:normal; letter-spacing:normal; line-height:20px; text-transform:none; text-align:right; text-decoration:underline; padding:0; margin:0">
									Oh wait, there's more! →
								</a>
							</td>
						</tr>
					</tbody></table>
					<!-- Content Table Close // -->
				</td>
			</tr>
		</tbody></table>
		<!-- Email Wrapper Webview Close //-->

		<!-- Email Wrapper Header Open //-->
		<table border="0" cellpadding="0" cellspacing="0" style="max-width:600px;" width="100%" class="wrapperWebview">
			<tbody><tr>
				<td align="center" valign="top">
					<!-- Content Table Open // -->
					<table border="0" cellpadding="0" cellspacing="0" width="100%">
						<tbody><tr>
							<td align="center" valign="middle" style="padding-top: 40px; padding-bottom: 40px;" class="emailLogo">
								<!-- Logo and Link // -->
								<a href="https://computing.president.ac.id" target="_blank" style="text-decoration:none;" class="">
									<img src="https://sg.pufacomputing.live/Logo%20Puma.png" alt="" width="150" border="0" style="width:100%; max-width:150px;height:auto; display:block;" class="">
								</a>
							</td>
						</tr>
					</tbody></table>
					<!-- Content Table Close // -->
				</td>
			</tr>
		</tbody></table>
		<!-- Email Wrapper Header Close //-->

		<!-- Email Wrapper Body Open // -->
		<table border="0" cellpadding="0" cellspacing="0" style="max-width:600px;" width="100%" class="wrapperBody">
			<tbody><tr>
				<td align="center" valign="top">

					<!-- Table Card Open // -->
					<table border="0" cellpadding="0" cellspacing="0" style="background-color:#FFFFFF;border-color:#E5E5E5; border-style:solid; border-width:0 1px 1px 1px;" width="100%" class="tableCard">

						<tbody><tr>
							<!-- Header Top Border // -->
							<td height="3" style="background-color:#003CE5;font-size:1px;line-height:3px;" class="topBorder">&nbsp;</td>
						</tr>


						<tr>
							<td align="center" valign="top" style="padding-bottom: 20px;" class="imgHero">
								<!-- Hero Image // -->
								<a href="https://computing.president.ac.id" target="_blank" style="text-decoration:none;">
									<img src="http://grapestheme.com/notify/img/hero-img/blue/heroFill/user-code.png" width="600" alt="" border="0" style="width:100%; max-width:600px; height:auto; display:block;">
								</a>
							</td>
						</tr>

						<tr>
							<td align="center" valign="top" style="padding-bottom: 5px; padding-left: 20px; padding-right: 20px;" class="mainTitle">
								<!-- Main Title Text // -->
								<h2 class="text" style="color:#000000; font-family:'Poppins', Helvetica, Arial, sans-serif; font-size:28px; font-weight:500; font-style:normal; letter-spacing:normal; line-height:36px; text-transform:none; text-align:center; padding:0; margin:0">
									Your OTP Code
								</h2>
							</td>
						</tr>

						<tr>
							<td align="center" valign="top" style="padding-bottom: 30px; padding-left: 20px; padding-right: 20px;" class="subTitle">
								<!-- Sub Title Text // -->
								<h4 class="text" style="color:#999999; font-family:'Poppins', Helvetica, Arial, sans-serif; font-size:16px; font-weight:500; font-style:normal; letter-spacing:normal; line-height:24px; text-transform:none; text-align:center; padding:0; margin:0">
									Unlock your account with this code
								</h4>
							</td>
						</tr>

						<tr>
							<td align="center" valign="top" style="padding-left:20px;padding-right:20px;" class="containtTable ui-sortable">

								<table border="0" cellpadding="0" cellspacing="0" width="100%" class="tableMediumTitle" style="">
									<tbody><tr>
										<td align="center" valign="top" style="padding-bottom: 20px;" class="mediumTitle">
											<!-- Medium Title Text // -->
											<p class="text" style="color:#3f4b97; font-family:'Poppins', Helvetica, Arial, sans-serif; font-size:34px; font-weight:300; font-style:normal; letter-spacing:normal; line-height:24px; text-transform:none; text-align:center; padding:0; margin:0">
												USE CODE : ` + data.OTPCode + `
											</p>
										</td>
									</tr>
								</tbody></table>

								<table border="0" cellpadding="0" cellspacing="0" width="100%" class="tableDescription" style="">
									<tbody><tr>
										<td align="center" valign="top" style="padding-bottom: 20px;" class="description">
											<!-- Description Text// -->
											<p class="text" style="color:#666666; font-family:'Open Sans', Helvetica, Arial, sans-serif; font-size:14px; font-weight:400; font-style:normal; letter-spacing:normal; line-height:22px; text-transform:none; text-align:center; padding:0; margin:0">
												Please use this code to unlock your account. This code will expire in 5 minutes.
											</p>
										</td>
									</tr>
								</tbody></table>

								

							</td>
						</tr>

						<tr>
							<td height="20" style="font-size:1px;line-height:1px;">&nbsp;</td>
						</tr>

						<tr><td height="20" style="font-size:1px;line-height:1px;">&nbsp;</td>
</tr>
					</tbody></table>
					<!-- Table Card Close// -->

					<!-- Space -->
					<table border="0" cellpadding="0" cellspacing="0" width="100%" class="space">
						<tbody><tr>
							<td height="30" style="font-size:1px;line-height:1px;">&nbsp;</td>
						</tr>
					</tbody></table>

				
                    
                    

                    
                    
</td>
			</tr>
		</tbody></table>
		<!-- Email Wrapper Body Close // -->

		<!-- Email Wrapper Footer Open // -->
		<table border="0" cellpadding="0" cellspacing="0" style="max-width:600px;" width="100%" class="wrapperFooter">
			<tbody><tr>
				<td align="center" valign="top">
					<!-- Content Table Open// -->
					<table border="0" cellpadding="0" cellspacing="0" width="100%" class="footer">
						<tbody><tr>
							<td align="center" valign="top" style="padding-top:10px;padding-bottom:10px;padding-left:10px;padding-right:10px;" class="socialLinks">
								<!-- Social Links (Facebook)// -->
								
								<!-- Social Links (Twitter)// -->
								
								<!-- Social Links (Pintrest)// -->
								
								<!-- Social Links (Instagram)// -->
								
								<!-- Social Links (Linkdin)// -->
								
							</td>
						</tr>

						<tr>
							<td align="center" valign="top" style="padding: 10px 10px 5px;" class="brandInfo">
								<!-- Brand Information // -->
								<p class="text" style="color:#777777; font-family:'Open Sans', Helvetica, Arial, sans-serif; font-size:12px; font-weight:400; font-style:normal; letter-spacing:normal; line-height:20px; text-transform:none; text-align:center; padding:0; margin:0;">©&nbsp;PUFA Computing. | President University | Cikarang, Indonesia.
								</p>
							</td>
						</tr>

						<tr>
							<td align="center" valign="top" style="padding: 0px 10px 20px;" class="footerLinks">
								<!-- Use Full Links (Privacy Policy)// -->
								<p class="text" style="color:#777777; font-family:'Open Sans', Helvetica, Arial, sans-serif; font-size:12px; font-weight:400; font-style:normal; letter-spacing:normal; line-height:20px; text-transform:none; text-align:center; padding:0; margin:0;"></p>
							</td>
						</tr>

						<tr>
							<td align="center" valign="top" style="padding: 0px 10px 10px;" class="footerEmailInfo">
								<!-- Information of NewsLetter (Subscribe Info)// -->
								<p class="text" style="color:#777777; font-family:'Open Sans', Helvetica, Arial, sans-serif; font-size:12px; font-weight:400; font-style:normal; letter-spacing:normal; line-height:20px; text-transform:none; text-align:center; padding:0; margin:0;">
								If you have any quetions please contact us <a href="mailto:pufa.computing@president.ac.id" style="color:#777777;text-decoration:underline;" target="_blank">pufa.computing@president.ac.id.</a>
								</p>
							</td>
						</tr>

						<tr>
							<td align="center" valign="top" style="padding-top:10px;padding-bottom:10px;padding-left:10px;padding-right:10px;" class="appLinks">
								<!-- App Links (Anroid)// -->
								
								<!-- App Links (IOs)// -->
								
							</td>
						</tr>

						<!-- Space -->
						<tr>
							<td height="30" style="font-size:1px;line-height:1px;">&nbsp;</td>
						</tr>
					</tbody></table>
					<!-- Content Table Close// -->
				</td>
			</tr>

			<!-- Space -->
			<tr>
				<td height="30" style="font-size:1px;line-height:1px;">&nbsp;</td>
			</tr>
		</tbody></table>
		<!-- Email Wrapper Footer Close // -->

		<!--[if (gte mso 9)|(IE)]></td></tr></table><![endif]-->
		</td>
	</tr>
</tbody></table>
	</center>
</body>
</html>`

	log.Println("Sending email to: ", to)

	// Make sure to modify the sendEmail function to support sending HTML content
	if err := ms.sendEmail(to, subject, body); err != nil {
		return err
	}

	return nil
}

func (ms *MailgunService) SendVerificationEmail(to, token string, userId uuid.UUID) error {
	// Get User's name from the database
	user, err := app.GetUserByID(userId)
	if err != nil {
		return err
	}

	subject := "Email Verification"

	// Load BaseURL from configs
	baseURL := configs.LoadConfig().BaseURL

	// Construct the verification link
	verificationLink := baseURL + "/auth/verify-email?token=" + token

	data := EmailData{
		Name:             user.FirstName,
		VerificationLink: verificationLink,
	}

	// Construct the HTML email body
	body := `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office">
<head>
	<!--[if (gte mso 9)|(IE)]>
	<xml>
		<o:OfficeDocumentSettings>
			<o:AllowPNG/>
			<o:PixelsPerInch>96</o:PixelsPerInch>
		</o:OfficeDocumentSettings>
	</xml>
	<![endif]-->
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	<meta name="viewport" content="width=device-width" />
	<meta http-equiv="X-UA-Compatible" content="IE=edge" />
	<title>Notify - Notify Email Newsletter</title>

	<!-- Google Fonts Link -->
	<link href="https://fonts.googleapis.com/css?family=Open+Sans:300,400,600,700" rel="stylesheet" />
	<link href="https://fonts.googleapis.com/css?family=Lora:400,700" rel="stylesheet" />
	<style type="text/css">

		/*------ Client-Specific Style ------ */
		@-ms-viewport{width:device-width;}
		table, td{mso-table-lspace:0pt; mso-table-rspace:0pt;}
		img{-ms-interpolation-mode:bicubic; border: 0;}
		p, a, li, td, blockquote{mso-line-height-rule:exactly;}
		p, a, li, td, body, table, blockquote{-ms-text-size-adjust:100%; -webkit-text-size-adjust:100%;}
		#outlook a{padding:0;}
		.ReadMsgBody{width:100%;} .ExternalClass{width:100%;}
		.ExternalClass,.ExternalClass div,.ExternalClass font,.ExternalClass p,.ExternalClass span,.ExternalClass td,img{line-height:100%;}

		/*------ Reset Style ------ */
		*{-webkit-text-size-adjust:none;-webkit-text-resize:100%;text-resize:100%;}
		table{border-spacing: 0 !important;}
		h1, h2, h3, h4, h5, h6, p{display:block; Margin:0; padding:0;}
		img, a img{border:0; height:auto; outline:none; text-decoration:none;}
		#bodyTable, #bodyCell{ margin:0; padding:0; width:100%;}
		body {height:100%; margin:0; padding:0; width:100%;}

		.appleLinks a {color: #c2c2c2 !important; text-decoration: none;}
        span.preheader { display: none !important; }

		/*------ Google Font Style ------ */
		[style*="Open Sans"] {font-family:'Open Sans', Helvetica, Arial, sans-serif !important;}
		[style*="Lora"] {font-family:'Lora', Georgia, Times, serif !important;}

		/*------ General Style ------ */
		.wrapperWebview, .wrapperBody, .wrapperFooter{width:100%; max-width:600px; Margin:0 auto;}

		/*------ Column Layout Style ------ */
		.tableCard {text-align:center; font-size:0;}
		
		/*------ Images Style ------ */
		.imgHero img{ width:600px; height:auto; }
		
	</style>

	<style type="text/css">
		/*------ Media Width 480 ------ */
		@media screen and (max-width:640px) {
			table[class="wrapperWebview"]{width:100% !important; }
			table[class="wrapperEmailBody"]{width:100% !important; }
			table[class="wrapperFooter"]{width:100% !important; }
			td[class="imgHero"] img{ width:100% !important;}
			.hideOnMobile {display:none !important; width:0; overflow:hidden;}
		}
	</style>

</head>

<body style="background-color:#F9F9F9;">
<center>
		<table border="0" cellpadding="0" cellspacing="0" width="100%" style="table-layout: fixed; background-color: rgb(255, 255, 255);" id="bodyTable">
	<tbody><tr>
		<td align="center" valign="top" style="padding-right:10px;padding-left:10px;" id="bodyCell">
		<!--[if (gte mso 9)|(IE)]><table align="center" border="0" cellspacing="0" cellpadding="0" style="width:600px;" width="600"><tr><td align="center" valign="top"><![endif]-->

		<!-- Email Wrapper Webview Open //-->
		<table border="0" cellpadding="0" cellspacing="0" style="max-width:600px;" width="100%" class="wrapperWebview">
			<tbody><tr>
				<td align="center" valign="top">
					<!-- Content Table Open // -->
					<table border="0" cellpadding="0" cellspacing="0" width="100%">
						<tbody><tr>
							<td align="right" valign="middle" style="padding-top: 20px; padding-right: 0px;" class="webview">
								<!-- Email View in Browser // -->
								<a class="text hideOnMobile" href="https://computing.president.ac.id" target="_blank" style="color: rgb(119, 119, 119); font-family: 'Open Sans', Helvetica, Arial, sans-serif; font-size: 12px; font-weight: 400; font-style: normal; letter-spacing: normal; line-height: 20px; text-transform: none; text-align: right; text-decoration: underline; padding: 0px; margin: 0px; display: none;">
									Oh wait, there's more! →
								</a>
							</td>
						</tr>
					</tbody></table>
					<!-- Content Table Close // -->
				</td>
			</tr>
		</tbody></table>
		<!-- Email Wrapper Webview Close //-->

		<!-- Email Wrapper Header Open //-->
		<table border="0" cellpadding="0" cellspacing="0" style="max-width:600px;" width="100%" class="wrapperWebview">
			<tbody><tr>
				<td align="center" valign="top">
					<!-- Content Table Open // -->
					<table border="0" cellpadding="0" cellspacing="0" width="100%">
						<tbody><tr>
							<td align="center" valign="middle" style="padding-top: 40px; padding-bottom: 40px;" class="emailLogo">
								<!-- Logo and Link // -->
								<a href="#" target="_blank" style="text-decoration:none;" class="">
									<img src="" alt="" width="150" border="0" style="width:100%; max-width:150px;height:auto; display:block;" class="">
								</a>
							</td>
						</tr>
					</tbody></table>
					<!-- Content Table Close // -->
				</td>
			</tr>
		</tbody></table>
		<!-- Email Wrapper Header Close //-->

		<!-- Email Wrapper Body Open // -->
		<table border="0" cellpadding="0" cellspacing="0" style="max-width:600px;" width="100%" class="wrapperBody">
			<tbody><tr>
				<td align="center" valign="top">

					<!-- Table Card Open // -->
					<table border="0" cellpadding="0" cellspacing="0" style="background-color:#FFFFFF;border-color:#E5E5E5; border-style:solid; border-width:0 1px 1px 1px;" width="100%" class="tableCard">

						<tbody><tr>
							<!-- Header Top Border // -->
							<td height="3" style="background-color: rgb(0, 0, 0); font-size: 1px; line-height: 3px;" class="topBorder">&nbsp;</td>
						</tr>


						<tr>
							<td align="center" valign="top" style="padding-bottom: 20px;" class="imgHero">
								<!-- Hero Image // -->
								<a href="#" target="_blank" style="text-decoration:none;" class="">
									<img src="http://grapestheme.com/notify/img/hero-img/blue/heroFill/user-subscribe.png" width="580" alt="" border="0" style="width: 100%; max-width: 580px; height: auto; display: block;" class="">
								</a>
							</td>
						</tr>

						<tr>
							<td align="center" valign="top" style="padding-bottom: 5px; padding-left: 20px; padding-right: 20px;" class="mainTitle">
								<!-- Main Title Text // -->
								<h2 class="text" style="color: rgb(23, 23, 23); font-family: 'Poppins', Helvetica, Arial, sans-serif; font-size: 28px; font-weight: 500; font-style: normal; letter-spacing: normal; line-height: 36px; text-transform: none; text-align: center; padding: 0px; margin: 0px;">
									Hi ` + data.Name + `!
								</h2>
							</td>
						</tr>

						<tr>
							<td align="center" valign="top" style="padding-bottom: 30px; padding-left: 20px; padding-right: 20px;" class="subTitle">
								<!-- Sub Title Text // -->
								<h4 class="text" style="color: rgb(23, 23, 23); font-family: 'Poppins', Helvetica, Arial, sans-serif; font-size: 16px; font-weight: 500; font-style: normal; letter-spacing: normal; line-height: 24px; text-transform: none; text-align: center; padding: 0px; margin: 0px;">
									Welcome to PUFA Computing Website 👻
								</h4>
							</td>
						</tr>

						<tr>
							<td align="center" valign="top" style="padding-left:20px;padding-right:20px;" class="containtTable ui-sortable">

								<table border="0" cellpadding="0" cellspacing="0" width="100%" class="tableDescription" style="">
									<tbody><tr>
										<td align="center" valign="top" style="padding-bottom: 20px;" class="description">
											<!-- Description Text// -->
											<p class="text" style="color: rgb(102, 102, 102); font-family: 'Open Sans', Helvetica, Arial, sans-serif; font-size: 14px; font-weight: 500; font-style: normal; letter-spacing: normal; line-height: 22px; text-transform: none; text-align: center; padding: 0px; margin: 0px;">Use the link  below to verify your email and start exploring our website.</p>
										</td>
									</tr>
								</tbody></table><table border="0" cellpadding="0" cellspacing="0" width="100%" class="tableButton" style="">
									<tbody><tr>
										<td align="center" valign="top" style="padding-top:20px;padding-bottom:20px;">

											<!-- Button Table // -->
											<table align="center" border="0" cellpadding="0" cellspacing="0">
												<tbody><tr>
													<td align="center" class="ctaButton" style="background-color: rgb(23, 23, 23); padding: 12px 35px; border-radius: 50px;">
														<!-- Button Link // -->
														<a class="text" href="` + data.VerificationLink + `" target="_blank" style="color: rgb(255, 255, 255); font-family: 'Open Sans', Helvetica, Arial, sans-serif; font-size: 14px; font-weight: 500; font-style: normal; letter-spacing: normal; line-height: 22px; text-transform: none; text-decoration: none; display: inline-block;">
Verify Now
														</a>
													</td>
												</tr>
											</tbody></table>

										</td>
									</tr>
								</tbody></table><table border="0" cellpadding="0" cellspacing="0" width="100%" class="tableDescription" style="">
									<tbody><tr>
										<td align="center" valign="top" style="padding-bottom: 20px;" class="description">
											<!-- Description Text// -->
											<p class="text" style="color: rgb(102, 102, 102); font-family: 'Open Sans', Helvetica, Arial, sans-serif; font-size: 14px; font-weight: 500; font-style: normal; letter-spacing: normal; line-height: 22px; text-transform: none; text-align: center; padding: 0px; margin: 0px;">Need Help?</p>
										</td>
									</tr>
								</tbody></table><table border="0" cellpadding="0" cellspacing="0" width="100%" class="tableDescription" style="">
									<tbody><tr>
										<td align="center" valign="top" style="padding-bottom: 20px;" class="description">
											<!-- Description Text// -->
											<p class="text" style="color: rgb(102, 102, 102); font-family: 'Open Sans', Helvetica, Arial, sans-serif; font-size: 14px; font-weight: 500; font-style: normal; letter-spacing: normal; line-height: 22px; text-transform: none; text-align: center; padding: 0px; margin: 0px;" data-font="active">Please send any feedback or bug reports
to developer@computizen.me</p>
										</td>
									</tr>
								</tbody></table>

								

							</td>
						</tr>

						<tr>
							<td height="20" style="font-size:1px;line-height:1px;">&nbsp;</td>
						</tr>

						<tr><td height="20" style="font-size:1px;line-height:1px;">&nbsp;</td>
</tr>
					</tbody></table>
					<!-- Table Card Close// -->

					<!-- Space -->
					<table border="0" cellpadding="0" cellspacing="0" width="100%" class="space">
						<tbody><tr>
							<td height="30" style="font-size:1px;line-height:1px;">&nbsp;</td>
						</tr>
					</tbody></table>

				</td>
			</tr>
		</tbody></table>
		<!-- Email Wrapper Body Close // -->

		<!-- Email Wrapper Footer Open // -->
		<table border="0" cellpadding="0" cellspacing="0" style="max-width:600px;" width="100%" class="wrapperFooter">
			<tbody><tr>
				<td align="center" valign="top">
					<!-- Content Table Open// -->
					<table border="0" cellpadding="0" cellspacing="0" width="100%" class="footer">
						<tbody><tr>
							<td align="center" valign="top" style="padding-top:10px;padding-bottom:10px;padding-left:10px;padding-right:10px;" class="socialLinks">
								<!-- Social Links (Facebook)// -->
								
								<!-- Social Links (Twitter)// -->
								
								<!-- Social Links (Pintrest)// -->
								
								<!-- Social Links (Instagram)// -->
								<a href="https://www.instagram.com/pucomputing/" target="_blank" style="display: inline-block;" class="instagram">
									<img src="http://grapestheme.com/notify/img/social/light/instagram.png" alt="" width="40" border="0" style="height:auto; width:100%; max-width:40px; margin-left:2px; margin-right:2px" class="">
								</a>
								<!-- Social Links (Linkdin)// -->
								<a href="https://www.linkedin.com/company/pufa-computing23/" target="_blank" style="display: inline-block;" class="linkdin">
									<img src="http://grapestheme.com/notify/img/social/light/linkdin.png" alt="" width="40" border="0" style="height:auto; width:100%; max-width:40px; margin-left:2px; margin-right:2px" class="">
								</a>
							<a href="https://www.youtube.com/@pufacomputingpresuniv" target="_blank" style="display: inline-block;" class="youtube"><img src="http://grapestheme.com/notify/img/social/light//youtube.png" alt="" width="40" border="0" style="height:auto; width:100%; max-width:40px; margin-left:2px; margin-right:2px" class=""></a></td>
						</tr>

						<tr>
							<td align="center" valign="top" style="padding: 10px 10px 5px;" class="brandInfo">
								<!-- Brand Information // -->
								<p class="text" style="color:#777777; font-family:'Open Sans', Helvetica, Arial, sans-serif; font-size:12px; font-weight:400; font-style:normal; letter-spacing:normal; line-height:20px; text-transform:none; text-align:center; padding:0; margin:0;">©&nbsp;PUFA Computing. All rights reserved. | President University, Jababeka, Indonesia.
								</p>
							</td>
						</tr>

						<tr>
							<td align="center" valign="top" style="padding: 0px 10px 20px;" class="footerLinks">
								<!-- Use Full Links (Privacy Policy)// -->
								<p class="text" style="color:#777777; font-family:'Open Sans', Helvetica, Arial, sans-serif; font-size:12px; font-weight:400; font-style:normal; letter-spacing:normal; line-height:20px; text-transform:none; text-align:center; padding:0; margin:0;">			</p>
							</td>
						</tr>

						<tr>
							<td align="center" valign="top" style="padding: 0px 10px 10px;" class="footerEmailInfo">
								<!-- Information of NewsLetter (Subscribe Info)// -->
								<p class="text" style="color:#777777; font-family:'Open Sans', Helvetica, Arial, sans-serif; font-size:12px; font-weight:400; font-style:normal; letter-spacing:normal; line-height:20px; text-transform:none; text-align:center; padding:0; margin:0;">						</p>
							</td>
						</tr>

						<tr>
							<td align="center" valign="top" style="padding-top:10px;padding-bottom:10px;padding-left:10px;padding-right:10px;" class="appLinks">
								<!-- App Links (Anroid)// -->
								
								<!-- App Links (IOs)// -->
								
							</td>
						</tr>

						<!-- Space -->
						<tr>
							<td height="30" style="font-size:1px;line-height:1px;">&nbsp;</td>
						</tr>
					</tbody></table>
					<!-- Content Table Close// -->
				</td>
			</tr>

			<!-- Space -->
			<tr>
				<td height="30" style="font-size:1px;line-height:1px;">&nbsp;</td>
			</tr>
		</tbody></table>
		<!-- Email Wrapper Footer Close // -->

		<!--[if (gte mso 9)|(IE)]></td></tr></table><![endif]-->
		</td>
	</tr>
</tbody></table>
	</center>
</body>
</html>`

	log.Println("Sending email to: ", to)

	// Make sure to modify the sendEmail function to support sending HTML content
	if err := ms.sendEmail(to, subject, body); err != nil {
		return err
	}

	log.Println()

	return nil
}

func (ms *MailgunService) sendEmail(toEmail, subject, body string) error {
	message := ms.mailgun.NewMessage(
		ms.senderEmail,
		subject,
		body,
		toEmail,
	)

	message.SetHtml(body)

	_, _, err := ms.mailgun.Send(context.Background(), message)
	if err != nil {
		return err
	}
	return nil
}
