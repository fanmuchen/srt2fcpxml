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

    // Create a blank title for the entire timeline
    blankTitle := Title.NewTitle("Blank Title", 0, subtitles.Duration().Seconds())

    for index, item := range subtitles.Items {
        // Create a new Title for each subtitle
        textStyleDef := Title.NewTextStyleDef(index + 1)
        text := Title.NewContent(index+1, func(lines []astisub.Line) string {
            var os []string
            for _, l := range lines {
                os = append(os, l.String())
            }
            return strings.Join(os, "\n")
        }(item.Lines))
        subtitle := Title.NewTitle(item.String(), item.StartAt.Seconds(), item.EndAt.Seconds()).SetTextStyleDef(textStyleDef).SetText(text)

        // Add subtitles as child elements to the blank title
        blankTitle.AddTitle(subtitle)
    }

    // Create spine and attach the blank title to it
    spine := Spine.NewSpine().AddTitle(blankTitle)
    seq := Sequence.NewSequence(subtitles.Duration().Seconds()).SetSpine(spine)
    project := Project.NewProject(projectName).SetSequence(seq)
    event := Event.NewEvent().SetProject(project)
    library := Library.NewLibrary(projectName).SetEvent(event)
    fcpxml.SetLibrary(library)

    return xml.MarshalIndent(fcpxml, "", "    ")
}
