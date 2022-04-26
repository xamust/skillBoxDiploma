package systemsProject

import (
	"github.com/sirupsen/logrus"
	"server/internal/app/models"
	"sort"
)

type SystemsProject struct {
	Logger           *logrus.Logger
	Config           *Config
	ParsingDataFiles map[string]string
}

//sms system..
func (s *SystemsProject) GetSMSData() ([][]models.SMSData, error) {
	//sms
	//init sms service
	sms := &SMSSystem{
		logger:   s.Logger,
		check:    &CheckData{Config: s.Config},
		fileName: s.ParsingDataFiles,
	}
	dataSMS, err := sms.ReadSMS()
	if err != nil {
		s.Logger.Errorf(err.Error())
		return nil, err
	}
	models.FullCountryNameSMS(dataSMS)

	// костыль с данными,ссылочный тип с указателями %)
	dataSMSDouble := make([]models.SMSData, len(dataSMS))
	copy(dataSMSDouble, dataSMS)
	sort.Slice(dataSMS, func(i, j int) bool {
		return dataSMS[i].Provider < dataSMS[j].Provider
	})
	sort.Slice(dataSMSDouble, func(i, j int) bool {
		return dataSMSDouble[i].Country < dataSMSDouble[j].Country
	})
	return [][]models.SMSData{dataSMS, dataSMSDouble}, nil
}

//voice system...
func (s *SystemsProject) getVoiceData() ([]models.VoiceCallData, error) {

	//init voice system...
	voice := &VoiceCallSystem{
		logger:   s.Logger,
		check:    &CheckData{Config: s.Config},
		fileName: s.ParsingDataFiles,
	}

	dataVoice, err := voice.ReadVoiceData()
	if err != nil {
		s.Logger.Errorf(err.Error())
		return nil, err
	}
	return dataVoice, nil
}

func (s *SystemsProject) GetResultData() (*models.ResultSetT, error) {
	/*
		type item struct {
			dataSMS       [][]models.SMSData
			dastaMMS      [][]models.MMSData
			dataVoiceCall []models.VoiceCallData
			dataEmail     map[string][][]models.EmailData
			dataBilling   models.BillingData
			dataSupport   []int
			dataIncidents []models.IncidentData
			err           error
		}
		dataS := make(chan item)

		go func() {
			var sms item
			sms.dataSMS, sms.err = s.GetSMSData()
			dataS <- sms
		}()
		sms := <-dataS
		close(dataS)
		if sms.err != nil {
			s.Logger.Error(sms.err)
			return nil, sms.err
		}
	*/
	sms, err := s.GetSMSData()
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}
	voice, err := s.getVoiceData()
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	return &models.ResultSetT{
		SMS:       sms,
		MMS:       nil,
		VoiceCall: voice,
		Email:     nil,
		Billing:   models.BillingData{},
		Support:   nil,
		Incidents: nil,
	}, nil
}
