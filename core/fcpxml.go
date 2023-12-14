package core

import (
	"encoding/xml"
	"srt2fcpxml/core/FcpXML"
	"srt2fcpxml/core/FcpXML/Library"
	"srt2fcpxml/core/FcpXML/Library/Event"
	"srt2fcpxml/core/FcpXML/Library/Event/Project"
	"srt2fcpxml/core/FcpXML/Library/Event/Project/Sequence"
	"srt2fcpxml/core/FcpXML/Library/Event/Project/Sequence/Spine"
	"srt2fcpxml/core/FcpXML/Library/Event/Project/Sequence/Spine/Gap"
	"srt2fcpxml/core/FcpXML/Library/Event/Project/Sequence/Spine/Gap/Title"
	"srt2fcpxml/core/FcpXML/Resources"
	"strings"

	"github.com/asticode/go-astisub"
)

func Srt2FcpXmlExport(projectName string, frameDuration interface{}, subtitles *astisub.Subtitles, width, height int) ([]byte, error) {
	fcpxml := FcpXML.New()
	res := Resources.NewResources()
	res.SetEffect(Resources.NewEffect())
	format := Resources.NewFormat().
		SetWidth(width).
		SetHeight(height).
		SetFrameRate(frameDuration).Render()
	res.SetFormat(format)
	fcpxml.SetResources(res)

	// Create a single gap for the entire timeline
	totalDuration := subtitles.Duration().Seconds()
	gap := Gap.NewGap(totalDuration)

	// Loop through subtitles and attach each subtitle to the gap
	for index, item := range subtitles.Items {
		textStyleDef := Title.NewTextStyleDef(index + 1)
		text := Title.NewContent(index+1, func(lines []astisub.Line) string {
			var os []string
			for _, l := range lines {
				os = append(os, l.String())
			}
			return strings.Join(os, "\n")
		}(item.Lines))
		title := Title.NewTitle(item.String(), item.StartAt.Seconds(), item.EndAt.Seconds()).SetTextStyleDef(textStyleDef).SetText(text)
		title.AddParam(Title.NewParams("Position", "9999/999166631/999166633/1/100/101", "0 -450"))
		title.AddParam(Title.NewParams("Alignment", "9999/999166631/999166633/2/354/999169573/401", "1 (Center)"))
		title.AddParam(Title.NewParams("Flatten", "9999/999166631/999166633/2/351", "1"))

		// Attach subtitle to the gap
		gap.AddTitle(title)
	}

	// Create a single spine for the entire timeline and set the gap
	spine := Spine.NewSpine().SetGap(gap)

	// Create a single sequence for the entire timeline and set the spine
	seq := Sequence.NewSequence(totalDuration).SetSpine(spine)

	// Create the project, event, and library
	project := Project.NewProject(projectName).SetSequence(seq)
	event := Event.NewEvent().SetProject(project)
	library := Library.NewLibrary(projectName).SetEvent(event)
	fcpxml.SetLibrary(library)

	return xml.MarshalIndent(fcpxml, "", "    ")
}